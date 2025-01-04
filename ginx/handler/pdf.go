package handler

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// PDFHandler serves a given PDF file for requests to its route
// https://github.com/cheikhsimsol/7zip/blob/main/main.go
func PDFHandler(filePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Open the PDF file
		file, err := os.Open(filePath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Could not open PDF file: %v", err)
			return
		}
		defer file.Close()

		// Set headers for PDF response
		c.Header("Content-Type", "application/pdf")

		// the file name should never be hard-coded.
		c.Header("Content-Disposition", "inline; filename=\"download.pdf\"")

		// Stream the file to the response
		_, err = io.Copy(c.Writer, file)
		if err != nil {
			c.String(http.StatusInternalServerError, "Could not serve PDF file: %v", err)
			return
		}

		c.Status(http.StatusOK)
	}
}
