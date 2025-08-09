package tts

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

// Handler holds the TTS client and configuration
type Handler struct {
	client   *Client
	audioDir string
}

// NewHandler creates a new TTS handler
func NewHandler(appID, accessKey, cluster, audioDir, baseURL string) *Handler {
	client := NewClient(appID, accessKey, cluster, baseURL)

	// Ensure audio directory exists
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		fmt.Printf("Warning: failed to create audio directory: %v\n", err)
	}

	return &Handler{
		client:   client,
		audioDir: audioDir,
	}
}

// SynthesizeHandler handles text-to-speech synthesis requests
func (h *Handler) SynthesizeHandler(c echo.Context) error {
	// Validate required fields
	text := c.QueryParam("text")
	if text == "" {
		return c.JSON(http.StatusBadRequest, TTSErrorResponse{
			Error: "Text is required",
			Code:  400,
		})
	}

	// Validate text length (max 2000 characters for most TTS services)
	if len(text) > 2000 {
		return c.JSON(http.StatusBadRequest, TTSErrorResponse{
			Error: "Text too long (max 2000 characters)",
			Code:  400,
		})
	}

	// Generate unique filename
	encodedFilename := base64.URLEncoding.EncodeToString([]byte(text))
	var shortFilename string
	if len(encodedFilename) >= 32 {
		shortFilename = encodedFilename[:32]
	} else {
		shortFilename = encodedFilename
	}
	filename := fmt.Sprintf("tts_%s.mp3", shortFilename)

	filepath := filepath.Join(h.audioDir, filename)

	// check if file exists
	if _, err := os.Stat(filepath); err == nil {
		// Set appropriate headers for audio streaming
		c.Response().Header().Set("Content-Type", "audio/mpeg")
		c.Response().Header().Set("Content-Disposition", "inline")
		// Allow seeking
		c.Response().Header().Set("Accept-Ranges", "bytes")
		return c.File(filepath)
	}

	// Call real TTS service
	audioData, err := h.client.SynthesizeText(text)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, TTSErrorResponse{
			Error:   "Failed to synthesize speech",
			Code:    500,
			Details: err.Error(),
		})
	}

	// Save audio file
	if err := os.WriteFile(filepath, audioData, 0644); err != nil {
		return c.JSON(http.StatusInternalServerError, TTSErrorResponse{
			Error:   "Failed to save audio file",
			Code:    500,
			Details: err.Error(),
		})
	}

	// Set appropriate headers for audio streaming
	c.Response().Header().Set("Content-Type", "audio/mpeg")
	c.Response().Header().Set("Content-Disposition", "inline")
	// Allow seeking
	c.Response().Header().Set("Accept-Ranges", "bytes")

	// Return response
	return c.File(filepath)
}

// ListAudioFilesHandler lists all audio files in the directory
func (h *Handler) ListAudioFilesHandler(c echo.Context) error {
	files, err := os.ReadDir(h.audioDir)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, TTSErrorResponse{
			Error:   "Failed to read audio directory",
			Code:    500,
			Details: err.Error(),
		})
	}

	var audioFiles []map[string]any
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Only include audio files
		ext := filepath.Ext(file.Name())
		if ext != ".mp3" && ext != ".wav" && ext != ".ogg" {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		audioFiles = append(audioFiles, map[string]any{
			"name":        file.Name(),
			"size":        info.Size(),
			"modified_at": info.ModTime().Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"files": audioFiles,
		"count": len(audioFiles),
	})
}
