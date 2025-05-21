package router

import (
	"divar_recommender/internal/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoutes(router *gin.Engine, webhookHandler *handlers.WebhookHandler) {
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Divar Recommender")
	})
	router.POST("/webhook", webhookHandler.HandleWebhook)
}
