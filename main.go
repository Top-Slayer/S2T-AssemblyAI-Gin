package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/gin-gonic/gin"
)

func main() {
	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	if apiKey == "" {
		log.Fatal("API key is not set in the environment variables.")
	}

	client := aai.NewClient(apiKey)
	ctx := context.Background()
	options := &aai.TranscriptOptionalParams{
		LanguageCode: "en",
		SpeechModel:  "best",
		BoostParam:   "high",
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/upload-audio", func(c *gin.Context) {
		file, err := c.FormFile("audio")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get audio file",
			})
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open the uploaded file"})
			return
		}
		defer openedFile.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, openedFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the file into memory"})
			return
		}

		transcript, _ := client.Transcripts.TranscribeFromReader(ctx, buf, options)

		c.JSON(http.StatusOK, gin.H{
			"message": *transcript.Text,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
