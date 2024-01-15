package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/quangkhaidam93/go-web-server/chat"
	"github.com/quangkhaidam93/go-web-server/controllers"
	"github.com/quangkhaidam93/go-web-server/example/kafka"
	"github.com/quangkhaidam93/go-web-server/middlewares"
	"github.com/quangkhaidam93/go-web-server/models"
	"google.golang.org/grpc"
)

func main() {
	// gRPC Server
	go func() {
		lis, err := net.Listen("tcp", ":3002")

		if err != nil {
			panic("[Error] failed to listen on port 3002 due to: " + err.Error())
		}

		grpcServer := grpc.NewServer()

		chatServer := chat.Server{}

		chat.RegisterChatServiceServer(grpcServer, &chatServer)

		if err := grpcServer.Serve(lis); err != nil {
			panic("[Error] failed to start gRPC server due to: " + err.Error())
		}
	}()

	// Http Server
	users := []models.BasicUser{
		{ID: 1, Name: "Emma"},
		{ID: 2, Name: "Bruno"},
		{ID: 3, Name: "Rick"},
		{ID: 4, Name: "Lena"},
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	models.ConnectDatabase()

	// Kafka producer
	producer, err := kafka.SetupProducer()

	if err != nil {
		fmt.Println(err)
	}

	defer producer.Close()

	// Kafka consumer
	store := &kafka.NotificationStore{
		Data: make(kafka.UserNotifications),
	}

	ctx, cancel := context.WithCancel(context.Background())

	go kafka.SetupConsumerGroup(ctx, store)

	defer cancel()

	// Router handler

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	public := router.Group("/api")

	public.GET("/posts", controllers.FindPosts)
	public.GET("/posts/:id", controllers.FindPost)
	public.POST("/login", controllers.Login)
	public.POST("/sign-up", controllers.SignUp)
	public.GET("/ping", controllers.CallGrpc)
	public.POST("/send", kafka.SendMessageHandler(producer, users))

	protected := router.Group("/api")
	protected.Use(middlewares.JwtAuthMiddleware)

	protected.POST("/posts", controllers.CreatePost)
	protected.PATCH("/posts/:id", controllers.UpdatePost)
	protected.DELETE("/posts/:id", controllers.DeletePost)

	if err := router.Run(":3001"); err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
