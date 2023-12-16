package middlewares

import (
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ricardoalcantara/web-notify/internal/controllers"
	"github.com/ricardoalcantara/web-notify/internal/models"
	"github.com/ricardoalcantara/web-notify/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	// secret := os.Getenv("JWT_SECRET")
	return func(c *gin.Context) {
		tokenType, authToken := getToken(c)

		if len(authToken) == 0 || (tokenType != "Basic" && tokenType != "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.ErrorResponse{Error: "Not authorized"})
			return
		}

		if tokenType == "Basic" {
			decoded, err := base64.StdEncoding.DecodeString(authToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.ErrorResponse{Error: utils.PrintError(err)})
				return
			}

			cred := strings.Split(string(decoded), ":")
			if len(cred) != 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.ErrorResponse{Error: "Not authorized"})
				return
			}

			client, err := models.LoginCheck(cred[0], cred[1])
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.ErrorResponse{Error: "Not authorized"})
				return
			}

			c.Set("x-id", strconv.Itoa(int(client.ID)))
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusNotImplemented, controllers.ErrorResponse{Error: "Authentication not implemented"})
			// authorized, err := token.IsAuthorized(authToken, secret)
			// if !authorized {
			// 	c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.ErrorResponse{Error: utils.PrintError(err)})
			// 	return
			// }

			// accessToken, err := token.ExtractToken(authToken, secret)
			// if err != nil {
			// 	c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.ErrorResponse{Error: utils.PrintError(err)})
			// 	return
			// }
			// claims, err := token.ExtractClaims(accessToken)
			// if err != nil {
			// 	c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.ErrorResponse{Error: utils.PrintError(err)})
			// 	return
			// }

			// c.Set("x-id", claims.RegisteredClaims.Subject)
			// c.Next()
		}
	}
}

func SseAuthMiddleware() gin.HandlerFunc {
	secret := os.Getenv("SSE_JWT_SECRET")
	return func(c *gin.Context) {

		var token string
		var exists bool
		if token, exists = c.GetQuery("token"); !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.ErrorResponse{Error: "token required"})
			return
		}

		var claims jwt.RegisteredClaims
		_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.ErrorResponse{Error: utils.PrintError(err)})
			return
		}

		c.Set("x-user-id", claims.Subject)
		c.Next()
	}
}

func getToken(c *gin.Context) (string, string) {
	authHeader := c.Request.Header.Get("Authorization")

	t := strings.Split(authHeader, " ")
	if len(t) == 2 {
		return t[0], t[1]
	}

	return "", ""
}
