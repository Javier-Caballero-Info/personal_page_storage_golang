package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/internal_services"
)

type FileController struct {
	fileService internal_services.FileService
}

func NewFileController(fileService internal_services.FileService) FileController {

	return FileController{
		fileService: fileService,
	}

}

func (fc *FileController) GetAll(c *gin.Context) {
	directory := c.Param("directory")

	filesList := fc.fileService.GetAllFiles(directory)

	c.JSON(200, gin.H{
		"directory": directory,
		"files":     filesList,
	})

	}
func (fc *FileController) Post(c *gin.Context) {

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

	result, err := fc.fileService.UploadFile(directory, handler.Filename, file)

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(201, result)

}
func (fc *FileController) Delete(c *gin.Context) {

	filePath := c.Param("filePath")

	err := fc.fileService.DeleteFile(filePath)

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Unable to delete " + filePath,
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(204, gin.H{})

}