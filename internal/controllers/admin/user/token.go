package user

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ricardoalcantara/web-notify/internal/utils"
)

func token(c *gin.Context) {
	userId := c.Param("userId")

	expiry, err := strconv.Atoi(utils.GetEnv("SSE_JWT_LIFESPAN", "5"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.PrintError(err)})
		return
	}

	jti, err := uuid.NewUUID()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.PrintError(err)})
		return
	}

	secret := os.Getenv("SSE_JWT_SECRET")
	exp := time.Now().Add(time.Second * time.Duration(expiry))
	claims := jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(exp),
		ID:        jti.String(),
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := unsignedToken.SignedString([]byte(secret))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.PrintError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
