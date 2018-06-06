package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/appleboy/gin-jwt"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/internal_services"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/external_services"
	"os"
	"time"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/controllers"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())

	s3Service := external_services.NewS3Service(os.Getenv("AWS_BASE_PATH"))

	fileService := internal_services.FileService{S3Service: s3Service, BasePath: os.Getenv("AWS_BASE_PATH")}

	fileController := controllers.NewFileController(fileService)

	jwtSecret := os.Getenv("JWT_SECRET")

	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "JavierCaballeroInfoStorage",
		Key:        []byte(jwtSecret),
		SigningAlgorithm: "HS384",
		Authorizator: func(user interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc: time.Now,
	}

	auth := r.Group("/")

	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/*directory", fileController.GetAll)

		auth.POST("/*directory", fileController.Post)

		auth.DELETE("/*filePath", fileController.Delete)
	}

	// Allow all CORS
	r.Use(cors.Default())

	port := "3000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	r.Run("0.0.0.0:" + port)
}
