package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/quangkhaidam93/go-web-server/controllers"
	"github.com/quangkhaidam93/go-web-server/middlewares"
	"github.com/quangkhaidam93/go-web-server/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	models.ConnectDatabase()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	public := router.Group("/api")

	public.GET("/posts", controllers.FindPosts)
	public.GET("/posts/:id", controllers.FindPost)
	public.POST("/login", controllers.Login)
	public.POST("/sign-up", controllers.SignUp)

	protected := router.Group("/api")
	protected.Use(middlewares.JwtAuthMiddleware)

	protected.POST("/posts", controllers.CreatePost)
	protected.PATCH("/posts/:id", controllers.UpdatePost)
	protected.DELETE("/posts/:id", controllers.DeletePost)

	if err := router.Run(":3001"); err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
