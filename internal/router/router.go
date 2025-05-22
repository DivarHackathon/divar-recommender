package router

import (
	"divar_recommender/internal/config"
	"divar_recommender/internal/handlers"
	"divar_recommender/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoutes(router *gin.Engine, webhookHandler *handlers.WebhookHandler) {
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Divar Recommender")
	})
	router.GET("/test", func(c *gin.Context) {
		d := services.NewDivarService(config.AppConfig.Divar.BaseURL, config.AppConfig.Divar.APIKey)
		r := services.NewRecommenderService(d, 2, 2, 0.5)

		recommendations, err := r.GetRecommendations("AauA0k5_")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"recommendations": recommendations,
		})
	})
	router.POST("/webhook", webhookHandler.HandleWebhook)
}
