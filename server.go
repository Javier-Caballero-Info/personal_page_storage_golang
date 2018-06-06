package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/internal_services"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/external_services"
	"os"
)


func main() {
	r := gin.Default()

	s3Service := external_services.NewS3Service(os.Getenv("AWS_BASE_PATH"))

	fileService := internal_services.FileService{S3Service: s3Service, BasePath: os.Getenv("AWS_BASE_PATH")}

	r.GET("/*directory", func(c *gin.Context) {
		directory := c.Param("directory")

		filesList := fileService.GetAllFiles(directory)

		c.JSON(200, gin.H{
			"directory": directory,
			"files":     filesList,
		})

	})

	r.POST("/*directory", func(c *gin.Context) {

		r := c.Request

		directory := c.Param("directory")

		r.ParseMultipartForm(32 << 20)

		file, handler, err := r.FormFile("upload")

		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		result, err := fileService.UploadFile(directory, handler.Filename, file)

		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(201, result)
	})


	r.DELETE("/*filePath", func(c *gin.Context) {

		filePath := c.Param("filePath")

		err := fileService.DeleteFile(filePath)

		if err != nil {
			c.JSON(404, gin.H{
				"message": "Unable to delete " + filePath,
				"detail": err.Error(),
			})
			return
		}

		c.JSON(204, gin.H{})

	})

	port := "3000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	r.Run("0.0.0.0:" + port)
}
