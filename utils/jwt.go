package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var secret = []byte(os.Getenv("SECRET"))

func CreateToken() (string, error) {
	token, err := generateJWT()

	if err != nil {

		return "", err
	}

	return token, nil
}

func extractToken(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")

	if split := strings.Split(token, " "); len(split) == 2 {
		return split[1]
	}

	return ""
}

func CheckToken(c *gin.Context) error {
	tokenString := extractToken(c)

	if err := verityJWT(tokenString); err != nil {
		return err
	}

	return nil
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["user"] = "user"

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verityJWT(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return err
	}

	return nil
}
