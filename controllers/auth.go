package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangkhaidam93/go-web-server/models"
	"github.com/quangkhaidam93/go-web-server/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthInfoInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SignUp(c *gin.Context) {
	var input AuthInfoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	user := models.User{Username: input.Username, Password: hashedPassword}

	result := models.DB.Create(&user)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})

}

func Login(c *gin.Context) {

	var input AuthInfoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	var user models.User

	if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	token, err := utils.CreateToken()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"token": token})
}
