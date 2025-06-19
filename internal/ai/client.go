package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"axon/internal/logger"
)

// Client represents an AI API client
type Client struct {
	openRouterKey string
	geminiKey     string
	httpClient    *http.Client
}

// NewClient creates a new AI client
func NewClient(openRouterKey, geminiKey string) *Client {
	return &Client{
		openRouterKey: openRouterKey,
		geminiKey:     geminiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Request represents an AI request
type Request struct {
	Prompt    string
	Model     string
	MaxTokens int
	Context   []string
}

// Response represents an AI response
type Response struct {
	Text  string
	Error error
}

// Generate generates content using the specified AI model
func (c *Client) Generate(req Request) (*Response, error) {
	logger.Info("Starting AI generation with model: %s", req.Model)
	logger.LogRequest(req)

	var resp *Response
	var err error

	if strings.HasPrefix(req.Model, "google/") {
		logger.Debug("Using Gemini provider")
		resp, err = c.generateGemini(req)
	} else {
		logger.Debug("Using OpenRouter provider")
		resp, err = c.generateOpenRouter(req)
	}

	if err != nil {
		logger.Error("AI generation failed: %v", err)
	} else {
		logger.Info("AI generation completed")
		logger.LogResponse(resp)
	}

	return resp, err
}

// generateOpenRouter generates content using OpenRouter API
func (c *Client) generateOpenRouter(req Request) (*Response, error) {
	logger.Debug("Starting OpenRouter request")

	if c.openRouterKey == "" {
		logger.Error("OpenRouter API key not configured")
		return &Response{Error: fmt.Errorf("OpenRouter API key not configured")}, nil
	}

	// Build messages from context and prompt
	messages := make([]map[string]string, 0)
	for _, ctx := range req.Context {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": ctx,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": req.Prompt,
	})

	payload := map[string]interface{}{
		"model":      req.Model,
		"messages":   messages,
		"max_tokens": req.MaxTokens,
		"stream":     false, // Explicitly disable streaming
	}

	logger.Debug("OpenRouter payload: %+v", payload)

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal OpenRouter request: %v", err)
		return &Response{Error: err}, nil
	}

	logger.Debug("OpenRouter request JSON: %s", string(data))

	logger.Debug("Creating OpenRouter HTTP request")
	request, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		logger.Error("Failed to create OpenRouter HTTP request: %v", err)
		return &Response{Error: err}, nil
	}

	request.Header.Set("Authorization", "Bearer "+c.openRouterKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("HTTP-Referer", "https://github.com/axon-game")
	request.Header.Set("X-Title", "Axon Game")

	logger.Debug("Sending OpenRouter HTTP request")
	resp, err := c.httpClient.Do(request)
	if err != nil {
		logger.Error("OpenRouter HTTP request failed: %v", err)
		return &Response{Error: err}, nil
	}

	logger.Debug("OpenRouter response status: %s", resp.Status)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read OpenRouter response: %v", err)
		return &Response{Error: err}, nil
	}

	logger.Debug("OpenRouter response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		logger.Error("OpenRouter API error: %s - %s", resp.Status, string(body))
		return &Response{Error: fmt.Errorf("API error: %s", string(body))}, nil
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("Failed to parse OpenRouter response: %v", err)
		return &Response{Error: err}, nil
	}

	logger.Debug("Parsed OpenRouter response: %+v", result)

	if len(result.Choices) == 0 {
		logger.Error("No choices in OpenRouter response")
		return &Response{Error: fmt.Errorf("no response from API")}, nil
	}

	response := &Response{Text: result.Choices[0].Message.Content}
	logger.Info("OpenRouter request completed successfully")
	logger.Debug("Extracted content: %s", response.Text)
	return response, nil
}

// generateGemini generates content using Gemini API (placeholder implementation)
func (c *Client) generateGemini(req Request) (*Response, error) {
	if c.geminiKey == "" {
		return &Response{Error: fmt.Errorf("gemini API key not configured")}, nil
	}

	// This is a simplified implementation - in practice you'd use the actual Gemini API
	return &Response{
		Text: "[Gemini response placeholder - implement actual Gemini API integration]",
	}, nil
}

// GetBestModel returns the best model for a specific task
func (c *Client) GetBestModel(task string) string {
	// Use fastest free models for reliable response times
	switch task {
	case "world_building":
		return "mistralai/mistral-7b-instruct:free" // Very fast free model
	case "storytelling":
		return "mistralai/mistral-7b-instruct:free"
	case "rule_setting":
		return "mistralai/mistral-7b-instruct:free"
	case "dialog":
		return "mistralai/mistral-7b-instruct:free"
	default:
		return "mistralai/mistral-7b-instruct:free"
	}
}
