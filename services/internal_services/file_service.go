package internal_services

import (
	"github.com/Javier-Caballero-Info/personal_page_storage_golang/services/external_services"
	"io"
	"strings"
)

type FileService struct {

	BasePath string

	S3Service external_services.S3Service
}

func (fileService *FileService) GetAllFiles(directory string) []map[string]string {

	var filesList []map[string]string

	if len(directory) > 1 {
		if string(directory[len(directory)-1]) != "/" {
			directory += "/"
		}
	} else {
		directory = "/"
	}

	directory = fileService.BasePath + directory

	files, err := fileService.S3Service.ListFiles(directory)

	if err == nil {

		for _, item := range files {

			fileName := *item.Key

			sl := strings.Split(fileName, "/")

			name := sl[len(sl)-1]
			if name != "" {
				if len(directory) > 2 {
					if len(strings.Split(fileName[len(directory)-1:], "/")) < 3  && fileName != directory {

						item := map[string]string{
							"name": name,
							"path": fileName,
							"url": fileService.S3Service.GetFileUrl(fileName),
						}

						filesList = append(filesList, item)
					}
				}else{
					if len(strings.Split(fileName, "/")) < 2 {

						item := map[string]string{
							"name": name,
							"path": fileName,
							"url": fileService.S3Service.GetFileUrl(fileName),
						}

						filesList = append(filesList, item)
					}
				}
			}

		}

	}

	if filesList == nil {
		filesList = []map[string]string{}
	}

	return filesList

}

func (fileService FileService) UploadFile(directory string, filename string, file io.Reader) (map[string]string, error){

	if len(directory) > 1 {
		if string(directory[len(directory)-1]) != "/" {
			directory += "/"
		}
	} else {
		directory = "/"
	}

	directory = fileService.BasePath + directory

	err := fileService.S3Service.UploadFile(directory + "/" + filename, file)

	var item map[string]string

	if err == nil {

		item = map[string]string{
			"name": filename,
			"path": directory + "/" + filename,
			"url":  fileService.S3Service.GetFileUrl(directory + filename),
		}

	}

	return item, err

}

func (fileService FileService) DeleteFile(filePath string) error {

	return fileService.S3Service.DeleteFile(filePath)

}