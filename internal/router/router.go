package router

import (
	"divar_recommender/internal/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler) {
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Divar Recommender")
	})
	router.POST("/chat-webhook", chatHandler.HandleChatWebhook)
}
