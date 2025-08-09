package main

import (
	"fmt"
	"net/http"

	"github.com/chaosflutter/common-go/tts"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.CORS())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Echo Web App!",
			"status":  "running",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// API group
	api := e.Group("/capi")

	// Load TTS configuration
	ttsConfig := tts.LoadConfig()

	// Setup TTS routes
	if ttsConfig.IsValid() {
		ttsHandler := tts.NewHandler(
			ttsConfig.AppID,
			ttsConfig.AccessKey,
			ttsConfig.Cluster,
			ttsConfig.AudioDir,
			ttsConfig.BaseURL,
		)

		// TTS API group
		ttsAPI := api.Group("/tts")

		// TTS routes
		ttsAPI.GET("", ttsHandler.SynthesizeHandler)
		ttsAPI.GET("/files", ttsHandler.ListAudioFilesHandler)

		fmt.Println("✅ TTS service enabled")
		fmt.Printf("📂 Audio files will be saved to: %s\n", ttsConfig.AudioDir)
	} else {
		// Add a disabled TTS info endpoint
		api.GET("/tts/status", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]any{
				"status":  "disabled",
				"message": "TTS service is disabled. Please configure environment variables.",
				"required_env": []string{
					"VOLCENGINE_TTS_APP_ID",
					"VOLCENGINE_TTS_ACCESS_KEY",
					"VOLCENGINE_TTS_SECRET_KEY",
					"TTS_ENABLED=true",
				},
			})
		})
		fmt.Println("⚠️  TTS service disabled - missing configuration")
	}

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
