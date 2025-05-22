package handlers

import (
	"divar_recommender/internal/config"
	"divar_recommender/internal/services"
	"divar_recommender/internal/types"
	"fmt"
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
	fmt.Println("\n here")

	var payload types.WebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	conversationID := payload.NewChatbotMessage.Conversation.ID

	text := payload.NewChatbotMessage.Text
	words := strings.Split(strings.TrimSpace(text), " ")

	var textMsg types.ChatMessage

	if payload.Type != "NEW_CHATBOT_MESSAGE" || len(words) != 1 {
		textMsg = h.chatService.BuildOnlyText("Invalid Input")
		err := h.chatService.SendMessage(conversationID, textMsg)
		if err != nil {
			log.Println(err)
		}
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	textMsg = h.chatService.BuildOnlyText("Invalid Input")
	err := h.chatService.SendMessage(conversationID, textMsg)
	if err != nil {
		log.Println(err)
	}

	token := words[0]

	d := services.NewDivarService(config.AppConfig.Divar.BaseURL, config.AppConfig.Divar.APIKey)
	r := services.NewRecommenderService(d, config.AppConfig.Recommendation.ProductionYearHigh, config.AppConfig.Recommendation.ProductionYearLow, config.AppConfig.Recommendation.UsageCoefficient)

	posts, _ := r.GetRecommendations(token)
	ads := r.MapPostsToRecommendationPosts(posts)

	for _, ad := range ads {
		textMsg := h.chatService.BuildAdText(types.Ad(ad))
		err := h.chatService.SendMessage(conversationID, textMsg)
		if err != nil {
			log.Println(err)
		}
	}
}
