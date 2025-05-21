package router

import (
	"divar_recommender/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/chat-webhook", chatHandler.HandleChatWebhook)
	}
}
