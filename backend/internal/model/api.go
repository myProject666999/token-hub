package model

import "time"

type ChatMessage struct {
	Role    string      `json:"role" binding:"required"`
	Content interface{} `json:"content" binding:"required"`
	Name    string      `json:"name,omitempty"`
}

type ChatCompletionRequest struct {
	Model            string         `json:"model" binding:"required"`
	Messages         []ChatMessage  `json:"messages" binding:"required"`
	Temperature      float64        `json:"temperature,omitempty"`
	TopP             float64        `json:"top_p,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Stop             interface{}    `json:"stop,omitempty"`
	PresencePenalty  float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64        `json:"frequency_penalty,omitempty"`
	User             string         `json:"user,omitempty"`
	Tools            interface{}    `json:"tools,omitempty"`
	ToolChoice       interface{}    `json:"tool_choice,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type CompletionRequest struct {
	Model            string  `json:"model" binding:"required"`
	Prompt           string  `json:"prompt" binding:"required"`
	Temperature      float64 `json:"temperature,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
	Stream           bool    `json:"stream,omitempty"`
	Stop             string  `json:"stop,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
}

type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type EmbeddingRequest struct {
	Model string `json:"model" binding:"required"`
	Input string `json:"input" binding:"required"`
	User  string `json:"user,omitempty"`
}

type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type ModelListResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID          string `json:"id"`
		Object      string `json:"object"`
		Created     int64  `json:"created"`
		OwnedBy     string `json:"owned_by"`
		Provider    string `json:"provider"`
		Description string `json:"description"`
	} `json:"data"`
}

type APIKey struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UserID    uint   `json:"user_id" gorm:"not null;index"`
	Key       string `json:"key" gorm:"uniqueIndex;size:64;not null"`
	Name      string `json:"name" gorm:"size:100"`
	Status    int    `json:"status" gorm:"default:1"`
	LastUsed  *time.Time `json:"last_used"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type APIKeyCreateRequest struct {
	Name string `json:"name" binding:"required"`
}
