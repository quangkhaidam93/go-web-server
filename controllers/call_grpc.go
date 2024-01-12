package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangkhaidam93/go-web-server/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CallGrpc(c *gin.Context) {
	conn, err := grpc.Dial(":3002", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	defer conn.Close()

	client := chat.NewChatServiceClient(conn)

	response, err := client.SayHello(context.Background(), &chat.Message{Body: "Hello From Client!"})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Response from gRPC server: %v", response.Body)})
}
