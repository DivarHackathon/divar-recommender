package handlers

import (
	"divar_recommender/internal/config"
	"divar_recommender/internal/services"
	"divar_recommender/internal/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type WebhookHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(chatService *services.ChatService) *WebhookHandler {
	return &WebhookHandler{
		chatService: chatService,
	}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	log.Println("=== Starting HandleWebhook ===")

	var payload types.WebhookPayload
	log.Println("Initialized empty payload struct")

	log.Println("Attempting to bind JSON payload...")
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("Failed to bind JSON payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		log.Println("Returned bad request response")
		return
	}
	log.Println("Successfully bound JSON payload")

	conversationID := payload.NewChatbotMessage.Conversation.ID
	log.Printf("Extracted conversation ID: %s", conversationID)

	text := payload.NewChatbotMessage.Text
	log.Printf("Extracted message text: %s", text)

	words := strings.Split(strings.TrimSpace(text), " ")
	log.Printf("Split text into %d words: %v", len(words), words)

	var textMsg types.ChatMessage
	log.Println("Initialized textMsg variable")

	log.Printf("Validating payload - Type: %s, First word: %s, Words count: %d",
		payload.Type,
		func() string {
			if len(words) > 0 {
				return words[0]
			} else {
				return "EMPTY"
			}
		}(),
		len(words))

	if payload.Type != "NEW_CHATBOT_MESSAGE" || words[0] != "start" || len(words) != 2 {
		log.Println("Validation failed - sending invalid input message")

		textMsg = h.chatService.BuildOnlyText("Invalid Input")
		log.Println("Built invalid input message")

		log.Printf("Sending invalid input message to conversation: %s", conversationID)
		err := h.chatService.SendMessage(conversationID, textMsg)
		if err != nil {
			log.Printf("Error sending invalid input message: %v", err)
		} else {
			log.Println("Successfully sent invalid input message")
		}

		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		log.Println("Returned ignored status response")
		return
	}
	log.Println("Validation passed - proceeding with processing")

	// Note: This seems like duplicate code from the validation block above
	log.Println("Building another invalid input message (this might be unintended)")
	textMsg = h.chatService.BuildOnlyText("Invalid Input")
	log.Println("Built second invalid input message")

	log.Printf("Sending second message to conversation: %s", conversationID)
	err := h.chatService.SendMessage(conversationID, textMsg)
	if err != nil {
		log.Printf("Error sending second message: %v", err)
	} else {
		log.Println("Successfully sent second message")
	}

	token := words[1]
	log.Printf("Extracted token from second word: %s", token)

	log.Println("Initializing Divar service...")
	log.Printf("Divar config - BaseURL: %s, APIKey: %s",
		config.AppConfig.Divar.BaseURL,
		func() string {
			if len(config.AppConfig.Divar.APIKey) > 4 {
				return config.AppConfig.Divar.APIKey[:4] + "****"
			}
			return "****"
		}())
	d := services.NewDivarService(config.AppConfig.Divar.BaseURL, config.AppConfig.Divar.APIKey)
	log.Println("Divar service initialized")

	log.Println("Initializing Recommender service...")
	log.Printf("Recommender config - ProductionYearHigh: %v, ProductionYearLow: %v, UsageCoefficient: %v",
		config.AppConfig.Recommendation.ProductionYearHigh,
		config.AppConfig.Recommendation.ProductionYearLow,
		config.AppConfig.Recommendation.UsageCoefficient)
	r := services.NewRecommenderService(d, config.AppConfig.Recommendation.ProductionYearHigh, config.AppConfig.Recommendation.ProductionYearLow, config.AppConfig.Recommendation.UsageCoefficient)
	log.Println("Recommender service initialized")

	log.Printf("Getting recommendations for token: %s", token)
	posts, err := r.GetRecommendations(token)
	if err != nil {
		log.Printf("Error getting recommendations: %v", err)
	} else {
		log.Printf("Successfully retrieved %d posts", len(posts))
	}

	log.Println("Mapping posts to recommendation posts...")
	ads := r.MapPostsToRecommendationPosts(posts)
	log.Printf("Mapped to %d ads", len(ads))

	log.Println("Starting to send ads to conversation...")
	for i, ad := range ads {
		log.Printf("Processing ad %d/%d", i+1, len(ads))

		textMsg := h.chatService.BuildAdText(types.Ad(ad))
		log.Printf("Built ad text message for ad %d", i+1)

		log.Printf("Sending ad %d to conversation: %s", i+1, conversationID)
		err := h.chatService.SendMessage(conversationID, textMsg)
		if err != nil {
			log.Printf("Error sending ad %d: %v", i+1, err)
		} else {
			log.Printf("Successfully sent ad %d", i+1)
		}
	}

	log.Printf("Finished processing all ads. Total ads sent: %d", len(ads))
	log.Println("=== HandleWebhook completed ===")
}
