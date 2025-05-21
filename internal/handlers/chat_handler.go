package handlers

import (
	"divar_recommender/internal/services"
	"divar_recommender/internal/types"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) HandleChatWebhook(c *gin.Context) {
	fmt.Println("\n here")

	var payload struct {
		Type              string `json:"type"`
		NewChatbotMessage struct {
			Text         string `json:"text"`
			Conversation struct {
				ID string `json:"id"`
			} `json:"conversation"`
			Sender struct {
				Type string `json:"type"`
			} `json:"sender"`
			Type string `json:"type"`
		} `json:"new_chatbot_message"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if payload.Type != "NEW_CHATBOT_MESSAGE" || strings.ToLower(payload.NewChatbotMessage.Text) != "/start" {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	conversationID := payload.NewChatbotMessage.Conversation.ID

	ads := []types.Ad{
		{
			Title: "پژو پارس XU7P مدل ۱۴۰۲",
			Price: 730000000,
			Image: "https://s100.divarcdn.com/static/photo/afra/post/sample1.jpg",
			Token: "Aat4oEs8",
		},
		{
			Title: "پرشیا ELX سفید مدل ۱۴۰۱",
			Price: 680000000,
			Image: "https://s100.divarcdn.com/static/photo/afra/post/sample2.jpg",
			Token: "AatgXZYu",
		},
		{
			Title: "پارس صفر موتور جدید",
			Price: 799000000,
			Image: "https://s100.divarcdn.com/static/photo/afra/post/sample3.jpg",
			Token: "Aak0FUDy",
		},
	}

	for _, ad := range ads {
		imageMsg := h.chatService.BuildImagePreview(ad)
		err := h.chatService.SendMessage(conversationID, imageMsg)
		if err != nil {
			log.Println(err)
		}

		textMsg := h.chatService.BuildTextOnly(ad)
		err = h.chatService.SendMessage(conversationID, textMsg)
		if err != nil {
			log.Println(err)
		}

		linkMsg := h.chatService.BuildLinkButton(ad)
		err = h.chatService.SendMessage(conversationID, linkMsg)
		if err != nil {
			log.Println(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "3 ads sent"})
}
