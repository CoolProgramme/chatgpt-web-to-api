package main

import (
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/api/imitate"
	_ "github.com/linweiyuan/go-chatgpt-api/env"
	"github.com/linweiyuan/go-chatgpt-api/middleware"
	"log"
	"os"
)

func init() {
	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)
}

//goland:noinspection SpellCheckingInspection
func main() {
	router := gin.Default()

	router.Use(middleware.CORS())
	router.Use(middleware.Authorization())

	setupChatGPTWebToAPI(router)
	setupAdminAPI(router)
	router.NoRoute(api.Proxy)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, api.ReadyHint)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server: " + err.Error())
	}
}

func setupChatGPTWebToAPI(router *gin.Engine) {
	group := router.Group("/v1")
	{
		group.POST("/chat/completions", imitate.CreateChatCompletions)
	}
}

func setupAdminAPI(router *gin.Engine) {
	group := router.Group("/admin")
	{
		group.POST("/user/token", imitate.AdminTokenGet)

		group.POST("/user/add", imitate.AdminUserAdd)

		group.POST("/token/add", imitate.AdminTokenAdd)

		group.GET("/token/count", imitate.AdminTokenCount)
	}
}
