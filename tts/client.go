package tts

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Volcengine TTS client
type Client struct {
	AppID     string
	AccessKey string
	Cluster   string
	BaseURL   string
}

// NewClient creates a new TTS client
func NewClient(appID, accessKey, cluster, baseURL string) *Client {
	return &Client{
		AppID:     appID,
		AccessKey: accessKey,
		Cluster:   cluster,
		BaseURL:   baseURL,
	}
}

// SynthesizeText converts text to speech and returns audio data
func (c *Client) SynthesizeText(text string) ([]byte, error) {
	// Generate request ID
	reqID := fmt.Sprintf("tts_%d", time.Now().UnixNano())

	// Prepare request payload
	request := TTSRequest{}
	request.App.Appid = c.AppID
	request.App.Token = c.AccessKey
	request.App.Cluster = c.Cluster

	request.User.UID = "default_user"

	request.Audio.VoiceType = "en_female_sarah_new_conversation_wvae_bigtts"
	request.Audio.Encoding = "mp3"
	request.Audio.SpeedRatio = 1.0

	request.Request.ReqID = reqID
	request.Request.Text = text
	request.Request.Operation = "query"

	// Convert to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Prepare HTTP request
	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer;"+c.AccessKey+"")

	// Make request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if response is JSON (error) or binary (audio)
	var ttsResp TTSResponse
	if err := json.Unmarshal(body, &ttsResp); err == nil && ttsResp.Code != 3000 {
		return nil, fmt.Errorf("TTS API error: %s (code: %d)", ttsResp.Message, ttsResp.Code)
	}

	audioData, err := base64.StdEncoding.DecodeString(ttsResp.Data)

	if err != nil {
		return nil, fmt.Errorf("failed to decode audio data: %w", err)
	}

	// If it's not a JSON error response, assume it's audio data
	return audioData, nil
}
