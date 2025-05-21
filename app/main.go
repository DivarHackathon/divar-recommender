package main

import (
	"divar_recommender/internal/config"
	"divar_recommender/internal/handlers"
	"divar_recommender/internal/router"
	"divar_recommender/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config.LoadConfig(".")
	fmt.Println("Server running on port:", config.AppConfig.Server.Port)
	chatService := services.NewChatService(config.AppConfig.Divar.APIKey)
	chatHandler := handlers.NewChatHandler(chatService)
	ginRouter := gin.Default()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(gin.Logger())
	router.SetupRoutes(ginRouter, chatHandler)
	strPort := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Println("Server running on port:", strPort)
	err := ginRouter.Run(strPort)
	if err != nil {
		log.Fatal(err)
		return
	}
}
