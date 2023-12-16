package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ricardoalcantara/web-notify/internal"
	admin_message "github.com/ricardoalcantara/web-notify/internal/controllers/admin/message"
	admin_user "github.com/ricardoalcantara/web-notify/internal/controllers/admin/user"
	"github.com/ricardoalcantara/web-notify/internal/controllers/frontend/sse"
	"github.com/ricardoalcantara/web-notify/internal/models"
	"github.com/ricardoalcantara/web-notify/internal/notification"
	"github.com/ricardoalcantara/web-notify/internal/setup"
	"github.com/ricardoalcantara/web-notify/internal/utils"
)

func init() {
	setup.Env()
	models.ConnectDataBase()
	internal.StartPing()
	notification.StartNotification()
}

func main() {
	go console()
	frontend()
}
func console() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(_ string) bool { return true }
	r.Use(cors.New(config))

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Online",
		})
	})

	admin_message.RegisterRoutes(r)
	admin_user.RegisterRoutes(r)

	host := utils.GetEnv("CONSOLE_HOST", "")
	port := utils.GetEnv("CONSOLE_PORT", "28585")
	r.Run(host + ":" + port)
}

func frontend() {
	r := gin.Default()

	if value, ok := os.LookupEnv("CORS_ORIGIN"); ok {
		config := cors.DefaultConfig()
		if value == "*" {
			config.AllowOriginFunc = func(_ string) bool { return true }
		} else {
			config.AllowOrigins = strings.Split(value, ",")
		}

		config.AddAllowHeaders("Authorization")

		r.Use(cors.New(config))
	}

	var files []string
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	r.LoadHTMLFiles(files...)
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "pages/home/index.html", nil)
	})

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Online",
		})
	})

	sse.RegisterRoutes(r)

	host := utils.GetEnv("HOST", "")
	port := utils.GetEnv("PORT", "18585")
	r.Run(host + ":" + port)
}
