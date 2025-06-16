package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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
	if strings.HasPrefix(req.Model, "google/") {
		return c.generateGemini(req)
	}
	return c.generateOpenRouter(req)
}

// generateOpenRouter generates content using OpenRouter API
func (c *Client) generateOpenRouter(req Request) (*Response, error) {
	if c.openRouterKey == "" {
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
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return &Response{Error: err}, nil
	}

	request, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return &Response{Error: err}, nil
	}

	request.Header.Set("Authorization", "Bearer "+c.openRouterKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("HTTP-Referer", "https://github.com/axon-game")
	request.Header.Set("X-Title", "Axon Game")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return &Response{Error: err}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{Error: err}, nil
	}

	if resp.StatusCode != http.StatusOK {
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
		return &Response{Error: err}, nil
	}

	if len(result.Choices) == 0 {
		return &Response{Error: fmt.Errorf("no response from API")}, nil
	}

	return &Response{Text: result.Choices[0].Message.Content}, nil
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
	switch task {
	case "world_building":
		return "anthropic/claude-3.5-sonnet"
	case "storytelling":
		return "openai/gpt-4o"
	case "rule_setting":
		return "openai/gpt-4o-mini"
	case "dialog":
		return "anthropic/claude-3-haiku"
	default:
		return "openai/gpt-4o-mini"
	}
}
