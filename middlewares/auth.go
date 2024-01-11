package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangkhaidam93/go-web-server/utils"
)

func JwtAuthMiddleware(c *gin.Context) {
	if err := utils.CheckToken(c); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
	}

	c.Next()
}
