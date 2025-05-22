package models

import "time"

type Post struct {
	Token    string `json:"token"`
	Category string `json:"category"`
	City     string `json:"city"`
	District string `json:"district"`
	Data     struct {
		Color string `json:"color"`
		Price struct {
			Value int    `json:"value"`
			Mode  string `json:"mode"`
		} `json:"price"`
		NewPrice                    int               `json:"new_price"`
		Longitude                   float64           `json:"longitude"`
		Usage                       int               `json:"usage"`
		MotorStatus                 string            `json:"motor_status"`
		Year                        string            `json:"year"`
		Description                 string            `json:"description"`
		BrandModel                  string            `json:"brand_model"`
		ThirdPartyInsuranceDeadline string            `json:"third_party_insurance_deadline"`
		Title                       string            `json:"title"`
		Gearbox                     string            `json:"gearbox"`
		Images                      []string          `json:"images"`
		BodyStatus                  string            `json:"body_status"`
		Latitude                    float64           `json:"latitude"`
		FuelType                    string            `json:"fuel_type"`
		BodyChassisStatus           BodyChassisStatus `json:"body_chassis_status"`
		Exchange                    bool              `json:"exchange"`
		ChassisStatus               string            `json:"chassis_status"`
	} `json:"data"`
	State            string    `json:"state"`
	FirstPublishedAt time.Time `json:"first_published_at"`
	ChatEnabled      bool      `json:"chat_enabled"`
	Score            float32   `json:"score"`
}

type PostItem struct {
	Token          string    `json:"token"`
	Category       string    `json:"category"`
	LastModifiedAt time.Time `json:"last_modified_at"`
	City           string    `json:"city"`
	Title          string    `json:"title"`
	Price          struct {
		Mode  string `json:"mode"`
		Value string `json:"value"`
	} `json:"price"`
	VehiclesFields struct {
		Usage string `json:"usage"`
	} `json:"vehicles_fields"`

	Score float32
}

type GetPostsRequestModel struct {
	Category string `json:"category"`
	City     string `json:"city"`
	Query    Query  `json:"query"`
}

type BodyChassisStatus struct {
	BackChassisStatus  string `json:"back_chassis_status"`
	BodyStatus         string `json:"body_status"`
	FrontChassisStatus string `json:"front_chassis_status"`
}

type Query struct {
	BrandModel     []string       `json:"brand_model"`
	ProductionYear ProductionYear `json:"production_year"`
	Usage          Usage          `json:"usage"`
}

type ProductionYear struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Usage struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type PostItems struct {
	Data []PostItem `json:"posts"`
}

type RecommendationPost struct {
	Title string `json:"title"`
	Price int    `json:"price"`
	Image string `json:"image"`
	Token string `json:"token"`
	Score string `json:"score"`
}
