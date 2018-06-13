package main

import (
	"github.com/gin-gonic/gin"
		"github.com/appleboy/gin-jwt"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/internal_services"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/external_services"
	"os"
	"time"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/controllers"
	"fmt"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())

	s3Service := external_services.NewS3Service(os.Getenv("AWS_BASE_PATH"))

	fileService := internal_services.FileService{S3Service: s3Service, BasePath: os.Getenv("AWS_BASE_PATH")}

	fileController := controllers.NewFileController(fileService)

	jwtSecret := os.Getenv("JWT_SECRET_KEY")

	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "JavierCaballeroInfoStorage",
		Key:        []byte(jwtSecret),
		SigningAlgorithm: os.Getenv("JWT_SIGN_ALGORITHM"),
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
	r.Use(CORSMiddleware())

	port := "3000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	r.Run("0.0.0.0:" + port)
}
