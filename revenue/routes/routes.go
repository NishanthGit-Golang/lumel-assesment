package routes

import (
	"net/http"
	"os"

	"revenue/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/uploadFile", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}
		tempFilePath := "./saveFile/" + file.Filename
		if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
			return
		}
		if err := services.LoadData(tempFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		os.Remove(tempFilePath)

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded and processed successfully!"})
	})

	return r
}
