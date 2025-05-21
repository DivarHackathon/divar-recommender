package services

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"divar_recommender/internal/models"
)

type RecommenderService struct {
	divarService       *DivarService
	productionYearHigh int
	productionYearLow  int
	usageCoefficient   float32
}

const (
	PriceAndUsageWeight = 0.2

	InsuranceWeightCoefficient = 0.05
	ProductionYearWeight       = 0.2

	GearboxWeight     = 0.05
	MotorStatusWeight = 0.2
	BodyStatusWeight  = 0.2
	ChassisWeight     = 0.1
)

const (
	IntactBodyStatus  = 1
	DamagedBodyStatus = 0

	HealthyMotorStatus = 1
	DamagedMotorStatus = 0

	AutoGearbox   = 1
	ManualGearbox = 0

	BothChassisHealthy = 1
	OneChassisHealthy  = 0.5
	BothChassisDamaged = 0
)

func NewRecommenderService(service *DivarService, productionYearHigh int, productionYearLow int, usageCoefficient float32) RecommenderService {
	return RecommenderService{
		divarService:       service,
		productionYearHigh: productionYearHigh,
		productionYearLow:  productionYearLow,
		usageCoefficient:   usageCoefficient,
	}
}

func safeNormalize(value, min, max int) float32 {
	if max == min {
		return 0
	}
	return float32(value-min) / float32(max-min)
}

func reverseSafeNormalize(value, min, max int) float32 {
	if max == min {
		return 0
	}
	return 1 - float32(value-min)/float32(max-min)
}

func (h *RecommenderService) GetRecommendations(token string) ([]models.Post, error) {
	p, err := h.divarService.GetPost(token)
	if err != nil {
		return nil, err
	}

	requestBody := h.NewGetPostsRequestModel(p)

	posts, err := h.divarService.GetPosts(requestBody)
	if err != nil {
		return nil, err
	}

	topPosts := h.sortByScore(posts, token)
	return topPosts, nil
}

func (h *RecommenderService) sortByScore(posts []models.PostItem, token string) []models.Post {
	firstStepResult := []models.PostItem{}

	var maxPrice int
	var minPrice int
	var maxUsage int
	var minUsage int

	for _, p := range posts {
		if p.Token == token {
			continue
		}

		if p.Price.Mode != "مقطوع" {
			continue
		}

		usage, _ := strconv.Atoi(p.VehiclesFields.Usage)
		price, _ := strconv.Atoi(p.Price.Value)

		if price > maxPrice {
			maxPrice = price
		}

		if minPrice == 0 || price < minPrice {
			minPrice = price
		}

		if usage > maxUsage {
			maxUsage = usage
		}

		if minUsage == 0 || usage < minUsage {
			minUsage = usage
		}
	}

	for _, p := range posts {
		if p.Token == token {
			continue
		}

		if p.Price.Mode != "مقطوع" {
			continue
		}

		usage, _ := strconv.Atoi(p.VehiclesFields.Usage)
		price, _ := strconv.Atoi(p.Price.Value)

		normPrice := (price - minPrice) / (maxPrice - minPrice)
		normUsage := (usage - minUsage) / (maxUsage - minUsage)

		score := (PriceAndUsageWeight * float32(normPrice)) + (PriceAndUsageWeight * float32(normUsage))

		p.Score = float32(score)

		firstStepResult = append(firstStepResult, p)
	}

	sort.Slice(firstStepResult, func(i, j int) bool {
		return firstStepResult[i].Score < firstStepResult[j].Score
	})

	if len(firstStepResult) > 10 {
		firstStepResult = firstStepResult[:10]
	} else {
		firstStepResult = firstStepResult[:]
	}

	var wg sync.WaitGroup

	topPosts := []models.Post{}
	for _, p := range firstStepResult {
		wg.Add(1)
		go func() {
			defer wg.Done()

			post, err := h.divarService.GetPost(p.Token)
			if err != nil {
				return
			}

			post.Score = p.Score
			topPosts = append(topPosts, post)
		}()

	}

	wg.Wait()

	var maxProductionYear int
	var minProductionYear int

	var maxInsuranceMonths int
	var minInsuranceMonths int

	topsResult := []models.Post{}

	for _, topPost := range topPosts {

		topPost.Score = 0
		usage := topPost.Data.Usage
		price := topPost.Data.Price.Value

		if price > maxPrice {
			maxPrice = price
		}

		if minPrice == 0 || price < minPrice {
			minPrice = price
		}

		if usage > maxUsage {
			maxUsage = usage
		}

		if minUsage == 0 || usage < minUsage {
			minUsage = usage
		}

		productionYear := ArabicToEnglishDigits(topPost.Data.Year)
		if maxProductionYear == 0 || maxProductionYear < productionYear {
			maxProductionYear = productionYear
		}

		if minProductionYear == 0 || minProductionYear > productionYear {
			minProductionYear = productionYear
		}

		insuranceDeadlineMonths, _ := strconv.Atoi(topPost.Data.ThirdPartyInsuranceDeadline)
		if maxInsuranceMonths == 0 || maxInsuranceMonths < insuranceDeadlineMonths {
			maxInsuranceMonths = insuranceDeadlineMonths
		}

		if minInsuranceMonths == 0 || minInsuranceMonths > insuranceDeadlineMonths {
			minInsuranceMonths = insuranceDeadlineMonths
		}

		normPrice := reverseSafeNormalize(price, minPrice, maxPrice)

		normUsage := reverseSafeNormalize(usage, minUsage, maxUsage)

		normProductionYear := safeNormalize(productionYear, minProductionYear, maxProductionYear)
		normInsuranceDeadline := safeNormalize(insuranceDeadlineMonths, minInsuranceMonths, maxInsuranceMonths)

		finalScore := PriceAndUsageWeight*float32(normPrice) +
			PriceAndUsageWeight*float32(normUsage) +
			ProductionYearWeight*float32(normProductionYear) +
			InsuranceWeightCoefficient*float32(normInsuranceDeadline)

		if topPost.Data.BodyStatus == "intact" {
			finalScore += BodyStatusWeight * float32(IntactBodyStatus)
		} else {
			finalScore += BodyStatusWeight * float32(DamagedBodyStatus)
		}

		if topPost.Data.Gearbox == "automatic" {
			finalScore += GearboxWeight * float32(AutoGearbox)
		} else {
			finalScore += GearboxWeight * float32(ManualGearbox)
		}

		if topPost.Data.BodyChassisStatus.BackChassisStatus == "healthy" && topPost.Data.BodyChassisStatus.FrontChassisStatus == "healthy" {
			finalScore += ChassisWeight * BothChassisHealthy
		} else {
			if topPost.Data.BodyChassisStatus.BackChassisStatus == "healthy" && topPost.Data.BodyChassisStatus.FrontChassisStatus != "healthy" {
				finalScore += ChassisWeight * OneChassisHealthy
			}
			if topPost.Data.BodyChassisStatus.BackChassisStatus != "healthy" && topPost.Data.BodyChassisStatus.FrontChassisStatus == "healthy" {
				finalScore += ChassisWeight * OneChassisHealthy
			} else {
				finalScore += ChassisWeight * BothChassisDamaged
			}
		}

		if topPost.Data.MotorStatus == "healthy" {
			finalScore += MotorStatusWeight * HealthyMotorStatus
		} else {
			finalScore += MotorStatusWeight * DamagedMotorStatus
		}

		topPost.Score = finalScore
		topsResult = append(topsResult, topPost)
	}

	sort.Slice(topsResult, func(i, j int) bool {
		return topsResult[i].Score > topsResult[j].Score
	})

	return topsResult
}

func (h *RecommenderService) NewGetPostsRequestModel(m models.Post) models.GetPostsRequestModel {
	year := ArabicToEnglishDigits(m.Data.Year)

	return models.GetPostsRequestModel{
		Category: m.Category,
		City:     m.City,
		Query: models.Query{
			BrandModel: []string{m.Data.BrandModel},
			ProductionYear: models.ProductionYear{
				Min: year + h.productionYearHigh,
				Max: year + h.productionYearLow,
			},
			Usage: models.Usage{
				Min: int(float32(m.Data.Usage) - (float32(m.Data.Usage) * h.usageCoefficient)),
				Max: int(float32(m.Data.Usage) + (float32(m.Data.Usage) * h.usageCoefficient)),
			},
		},
	}
}

func ArabicToEnglishDigits(input string) int {
	var builder strings.Builder

	for _, r := range input {
		switch {
		case r >= '٠' && r <= '٩':
			builder.WriteRune('0' + (r - '٠'))
		case r >= '۰' && r <= '۹':
			builder.WriteRune('0' + (r - '۰'))
		default:
			builder.WriteRune(r)
		}
	}

	numStr := builder.String()
	num, _ := strconv.Atoi(numStr)
	return num
}
