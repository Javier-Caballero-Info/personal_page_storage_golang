package main

import (
	"github.com/gin-gonic/gin"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"fmt"
	"strings"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)



func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	r := gin.Default()

	region := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})

	svc := s3.New(sess)

	uploader := s3manager.NewUploader(sess)

	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}

	bucket := os.Getenv("AWS_BUCKET")
	basePath := os.Getenv("AWS_BASE_PATH")


	r.GET("/*directory", func(c *gin.Context) {
		directory := c.Param("directory")

		filesList := []map[string]interface{}{}

		if len(directory) > 1 {

			directory = directory[1:]

			if string(directory[len(directory)-1]) != "/" {
				directory += "/"
			}
		} else {
			directory = ""
		}

		directory += basePath

		resp, err := svc.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucket),
			Prefix: aws.String(directory),
		})

		if err != nil {
			exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
		}

		for _, item := range resp.Contents {

			fileName := *item.Key

			sl := strings.Split(fileName, "/")

			name := sl[len(sl)-1]

			if len(directory) > 2 {
				if len(strings.Split(fileName[len(directory)-1:], "/")) < 3  && fileName != directory {

					item := map[string]interface{}{
						"name": name,
						"path": fileName,
						"url": fmt.Sprintf(
							"https://s3.amazonaws.com/%s/%s",
							bucket,
							strings.Replace(fileName, " ", "+", -1),
						),
					}

					filesList = append(filesList, item)
				}
			}else{
				if len(strings.Split(fileName, "/")) < 2 {

					item := map[string]interface{}{
						"name": name,
						"path": fileName,
						"url": fmt.Sprintf(
							"https://s3.amazonaws.com/%s/%s",
							bucket,
							strings.Replace(fileName, " ", "+", -1),
						),
					}

					filesList = append(filesList, item)
				}
			}

		}

		c.JSON(200, gin.H{
			"directory": "/" + directory,
			"files":     filesList,
		})

	})

	r.POST("/*directory", func(c *gin.Context) {
		directory := c.Param("directory")

		r := c.Request

		if len(directory) > 1 {

			directory = directory[1:]

			if string(directory[len(directory)-1]) != "/" {
				directory += "/"
			}
		} else {
			directory = ""
		}

		directory += basePath

		r.ParseMultipartForm(32 << 20)

		file, handler, err := r.FormFile("upload")
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		filename := directory + handler.Filename

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key: aws.String(filename),
			Body: file,
			CacheControl: aws.String("max-age=86400"),
			ACL: aws.String("public-read"),
		})

		if err != nil {
			c.JSON(500, gin.H{
				"message": "Unable to upload " + filename,
				"detail": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"message": "Successfully uploaded " + filename,
		})
	})

	r.DELETE("/*filePath", func(c *gin.Context) {

		filePath := c.Param("filePath")

		_, err := svc.GetObject(&s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(filePath)})

		if err != nil {
			c.JSON(404, gin.H{
				"message": "Unable to delete " + filePath,
				"detail": err.Error(),
			})
			return
		}

		_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(filePath)})

		if err != nil {
			c.JSON(500, gin.H{
				"message": "Unable to upload " + filePath,
				"detail": fmt.Sprintf("Unable to delete object %q from bucket %q, %v", filePath, bucket, err),
			})
			return
		}

		err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filePath),
		})

		c.JSON(200, gin.H{
			"message": "Successfully deleted " + filePath,
		})

	})

	r.Run(":3000") // listen and serve on 0.0.0.0:3000
}
