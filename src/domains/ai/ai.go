package domainAI

import "context"

type IAIUsecase interface {
	Chat(ctx context.Context, request ChatRequest) (response ChatResponse, err error)
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message             string        `json:"message"`
	ConversationHistory []ChatMessage `json:"conversation_history"`
	SystemPrompt        string        `json:"system_prompt,omitempty"`
}

type ChatResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
}
