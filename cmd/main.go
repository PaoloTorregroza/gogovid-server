package main

import (
	"go-streamer/internal/handlers"
	"go-streamer/internal/handlers/oauth2"
	"go-streamer/internal/repositorioes"
	"go-streamer/internal/utils"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Repo := repositorioes.NewS3Repo()
	dbRepo := repositorioes.NewDBRepo(os.Getenv("DATABASE_ENGINE"))

	r := gin.Default()

	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true

	r.Use(cors.New(corsConf))
	r.Use(func(c *gin.Context) {
		c.Set(utils.S3_REPO_CTX_KEY, s3Repo)
		c.Set(utils.DB_REPO_CTX_KEY, dbRepo)
		c.Set(utils.GOOGLE_OAUTH_CONF_CTX_KEY, oauth2.NewGoogleOauthConfig())
		c.Set(utils.GITHUB_OAUTH_CONF_CTX_KEY, oauth2.NewGithubOauthConfig())
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		repo := c.MustGet(utils.S3_REPO_CTX_KEY).(*repositorioes.S3Repo)
		repo.TestListObject()
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Video routes
	r.GET("/video/:fileId", handlers.ServeVideo)
	r.POST("/video", handlers.AuthMiddleware, handlers.UploadVideo)

	// Oauth2 Routes
	oauth2.Oauth2Handler(r)

	r.Run()
}
