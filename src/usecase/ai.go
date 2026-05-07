package usecase

import (
        "bytes"
        "context"
        "encoding/json"
        "fmt"
        "io"
        "net/http"

        "sanjaywa/config"
        domainAI "sanjaywa/domains/ai"
        pkgError "sanjaywa/pkg/error"
)

type serviceAI struct{}

func NewAIService() domainAI.IAIUsecase {
        return &serviceAI{}
}

type groqMessage struct {
        Role    string `json:"role"`
        Content string `json:"content"`
}

type groqRequest struct {
        Model     string        `json:"model"`
        MaxTokens int           `json:"max_tokens"`
        Messages  []groqMessage `json:"messages"`
}

type groqResponse struct {
        Choices []struct {
                Message struct {
                        Content string `json:"content"`
                } `json:"message"`
        } `json:"choices"`
        Error *struct {
                Message string `json:"message"`
        } `json:"error,omitempty"`
}

func (s *serviceAI) Chat(ctx context.Context, request domainAI.ChatRequest) (domainAI.ChatResponse, error) {
        if config.GroqAPIKey == "" {
                return domainAI.ChatResponse{}, pkgError.InternalServerError("GROQ_API_KEY is not configured")
        }

        systemPrompt := request.SystemPrompt
        if systemPrompt == "" {
                systemPrompt = "You are a helpful AI assistant integrated with a WhatsApp gateway. Be concise and helpful."
        }

        messages := []groqMessage{
                {Role: "system", Content: systemPrompt},
        }

        for _, msg := range request.ConversationHistory {
                messages = append(messages, groqMessage{
                        Role:    msg.Role,
                        Content: msg.Content,
                })
        }

        messages = append(messages, groqMessage{
                Role:    "user",
                Content: request.Message,
        })

        payload := groqRequest{
                Model:     config.GroqModel,
                MaxTokens: config.GroqMaxTokens,
                Messages:  messages,
        }

        body, err := json.Marshal(payload)
        if err != nil {
                return domainAI.ChatResponse{}, pkgError.InternalServerError(fmt.Sprintf("failed to encode request: %v", err))
        }

        req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(body))
        if err != nil {
                return domainAI.ChatResponse{}, pkgError.InternalServerError(fmt.Sprintf("failed to build request: %v", err))
        }
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+config.GroqAPIKey)

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
                return domainAI.ChatResponse{}, pkgError.InternalServerError(fmt.Sprintf("failed to call Groq API: %v", err))
        }
        defer resp.Body.Close()

        respBody, err := io.ReadAll(resp.Body)
        if err != nil {
                return domainAI.ChatResponse{}, pkgError.InternalServerError(fmt.Sprintf("failed to read response: %v", err))
        }

        var groqResp groqResponse
        if err := json.Unmarshal(respBody, &groqResp); err != nil {
                return domainAI.ChatResponse{}, pkgError.InternalServerError(fmt.Sprintf("failed to parse response: %v", err))
        }

        if groqResp.Error != nil {
                return domainAI.ChatResponse{}, pkgError.InternalServerError(fmt.Sprintf("Groq API error: %s", groqResp.Error.Message))
        }

        if len(groqResp.Choices) == 0 || groqResp.Choices[0].Message.Content == "" {
                return domainAI.ChatResponse{}, pkgError.InternalServerError("no response received from Groq API")
        }

        return domainAI.ChatResponse{
                Response: groqResp.Choices[0].Message.Content,
                Model:    config.GroqModel,
        }, nil
}
