package main

import (
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/controllers"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/external_services"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/internal_services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

// middleware to protect private pages
func ValidateJWT()  gin.HandlerFunc {
	return func (c *gin.Context) {
		// middleware
		jwtSecret := os.Getenv("JWT_SECRET_KEY")

		auth := c.GetHeader("Authorization")

		if auth != "" {
			tokens := strings.Fields(auth)
			if len(tokens) == 2 && tokens[0] == "Bearer"{
				token, err := jwt.Parse(tokens[1], func(token *jwt.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				})

				if err == nil && token.Valid {
					c.Next()
				} else {
					c.AbortWithStatusJSON(401, gin.H{
						"code":    401,
						"message": "Unauthorized",
					})
				}
			}
		} else {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "Unauthorized",
			})
		}

	}
}


func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Origin")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
	}
}

func main() {
	r := gin.Default()

	r.Use(gin.Logger())

	r.Use(CORSMiddleware())

	r.OPTIONS("/*path", CORSMiddleware())

	s3Service := external_services.NewS3Service(os.Getenv("AWS_BASE_PATH"))

	fileService := internal_services.FileService{S3Service: s3Service, BasePath: os.Getenv("AWS_BASE_PATH")}

	fileController := controllers.NewFileController(fileService)

	auth := r.Group("/")

	auth.Use(ValidateJWT())
	{
		auth.GET("/*directory", fileController.GetAll)

		auth.POST("/*directory", fileController.Post)

		auth.DELETE("/*filePath", fileController.Delete)
	}

	port := "3000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	r.Run("0.0.0.0:" + port)
}
