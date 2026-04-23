package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"token-hub/internal/model"
	"token-hub/internal/repository"

	"github.com/google/uuid"
)

type APIService struct {
	pointsService *PointsService
	callLogService *CallLogService
}

func NewAPIService() *APIService {
	return &APIService{
		pointsService:  NewPointsService(),
		callLogService: NewCallLogService(),
	}
}

func (s *APIService) ChatCompletion(userID uint, req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, error) {
	requestID := uuid.New().String()
	startTime := time.Now()

	modelInfo, err := s.getActiveModel(req.Model)
	if err != nil {
		s.callLogService.CreateLog(&model.CallLog{
			UserID:      userID,
			Model:       *modelInfo,
			RequestID:   requestID,
			Status:      model.CallStatusFailed,
			ErrorMessage: err.Error(),
			Duration:    time.Since(startTime).Milliseconds(),
		})
		return nil, err
	}

	userPoints, err := s.checkUserPoints(userID)
	if err != nil {
		s.callLogService.CreateLog(&model.CallLog{
			UserID:       userID,
			ProviderID:   modelInfo.ProviderID,
			ModelID:      modelInfo.ID,
			RequestID:    requestID,
			Status:       model.CallStatusFailed,
			ErrorMessage: err.Error(),
			Duration:     time.Since(startTime).Milliseconds(),
		})
		return nil, err
	}

	estimatedInputTokens := s.estimateInputTokens(req.Messages)
	estimatedCost := s.estimateCost(modelInfo, estimatedInputTokens, 0)

	if userPoints < estimatedCost {
		return nil, errors.New("积分不足，无法完成请求")
	}

	response, inputTokens, outputTokens, err := s.callModelAPI(modelInfo, req)
	if err != nil {
		s.callLogService.CreateLog(&model.CallLog{
			UserID:       userID,
			ProviderID:   modelInfo.ProviderID,
			ModelID:      modelInfo.ID,
			RequestID:    requestID,
			InputTokens:  inputTokens,
			OutputTokens: outputTokens,
			TotalTokens:  inputTokens + outputTokens,
			Status:       model.CallStatusFailed,
			ErrorMessage: err.Error(),
			Duration:     time.Since(startTime).Milliseconds(),
		})
		return nil, err
	}

	totalTokens := inputTokens + outputTokens
	pointsConsumed := s.calculatePoints(modelInfo, inputTokens, outputTokens)

	if err := s.pointsService.ConsumePoints(
		userID,
		pointsConsumed,
		fmt.Sprintf("调用模型 %s，输入%d tokens，输出%d tokens", modelInfo.Name, inputTokens, outputTokens),
		requestID,
		"api_call",
	); err != nil {
		s.callLogService.CreateLog(&model.CallLog{
			UserID:         userID,
			ProviderID:     modelInfo.ProviderID,
			ModelID:        modelInfo.ID,
			RequestID:      requestID,
			InputTokens:    inputTokens,
			OutputTokens:   outputTokens,
			TotalTokens:    totalTokens,
			PointsConsumed: pointsConsumed,
			Status:         model.CallStatusFailed,
			ErrorMessage:   "积分扣减失败: " + err.Error(),
			Duration:       time.Since(startTime).Milliseconds(),
		})
		return nil, errors.New("积分扣减失败: " + err.Error())
	}

	s.callLogService.CreateLog(&model.CallLog{
		UserID:         userID,
		ProviderID:     modelInfo.ProviderID,
		ModelID:        modelInfo.ID,
		RequestID:      requestID,
		InputTokens:    inputTokens,
		OutputTokens:   outputTokens,
		TotalTokens:    totalTokens,
		PointsConsumed: pointsConsumed,
		Status:         model.CallStatusSuccess,
		Duration:       time.Since(startTime).Milliseconds(),
	})

	return response, nil
}

func (s *APIService) Completion(userID uint, req *model.CompletionRequest) (*model.CompletionResponse, error) {
	requestID := uuid.New().String()
	startTime := time.Now()

	modelInfo, err := s.getActiveModel(req.Model)
	if err != nil {
		return nil, err
	}

	userPoints, err := s.checkUserPoints(userID)
	if err != nil {
		return nil, err
	}

	estimatedInputTokens := len(req.Prompt) / 4
	estimatedCost := s.estimateCost(modelInfo, estimatedInputTokens, 0)

	if userPoints < estimatedCost {
		return nil, errors.New("积分不足，无法完成请求")
	}

	chatReq := &model.ChatCompletionRequest{
		Model: req.Model,
		Messages: []model.ChatMessage{
			{Role: "user", Content: req.Prompt},
		},
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		MaxTokens:        req.MaxTokens,
		Stream:           req.Stream,
		Stop:             req.Stop,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
	}

	chatResp, inputTokens, outputTokens, err := s.callModelAPI(modelInfo, chatReq)
	if err != nil {
		s.callLogService.CreateLog(&model.CallLog{
			UserID:       userID,
			ProviderID:   modelInfo.ProviderID,
			ModelID:      modelInfo.ID,
			RequestID:    requestID,
			InputTokens:  inputTokens,
			OutputTokens: outputTokens,
			TotalTokens:  inputTokens + outputTokens,
			Status:       model.CallStatusFailed,
			ErrorMessage: err.Error(),
			Duration:     time.Since(startTime).Milliseconds(),
		})
		return nil, err
	}

	totalTokens := inputTokens + outputTokens
	pointsConsumed := s.calculatePoints(modelInfo, inputTokens, outputTokens)

	if err := s.pointsService.ConsumePoints(
		userID,
		pointsConsumed,
		fmt.Sprintf("调用模型 %s，输入%d tokens，输出%d tokens", modelInfo.Name, inputTokens, outputTokens),
		requestID,
		"api_call",
	); err != nil {
		return nil, errors.New("积分扣减失败: " + err.Error())
	}

	s.callLogService.CreateLog(&model.CallLog{
		UserID:         userID,
		ProviderID:     modelInfo.ProviderID,
		ModelID:        modelInfo.ID,
		RequestID:      requestID,
		InputTokens:    inputTokens,
		OutputTokens:   outputTokens,
		TotalTokens:    totalTokens,
		PointsConsumed: pointsConsumed,
		Status:         model.CallStatusSuccess,
		Duration:       time.Since(startTime).Milliseconds(),
	})

	resp := &model.CompletionResponse{
		ID:      chatResp.ID,
		Object:  "text_completion",
		Created: chatResp.Created,
		Model:   chatResp.Model,
	}

	for _, choice := range chatResp.Choices {
		resp.Choices = append(resp.Choices, struct {
			Text         string `json:"text"`
			Index        int    `json:"index"`
			FinishReason string `json:"finish_reason"`
		}{
			Text:         fmt.Sprintf("%v", choice.Message.Content),
			Index:        choice.Index,
			FinishReason: choice.FinishReason,
		})
	}

	resp.Usage = chatResp.Usage

	return resp, nil
}

func (s *APIService) GetAvailableModels(userID uint) (*model.ModelListResponse, error) {
	var providers []model.Provider
	if err := repository.DB.Preload("Models", "status = ?", 1).
		Where("status = ?", 1).
		Order("sort ASC").
		Find(&providers).Error; err != nil {
		return nil, err
	}

	response := &model.ModelListResponse{
		Object: "list",
	}

	for _, provider := range providers {
		for _, m := range provider.Models {
			response.Data = append(response.Data, struct {
				ID          string `json:"id"`
				Object      string `json:"object"`
				Created     int64  `json:"created"`
				OwnedBy     string `json:"owned_by"`
				Provider    string `json:"provider"`
				Description string `json:"description"`
			}{
				ID:          m.Code,
				Object:      "model",
				Created:     m.CreatedAt.Unix(),
				OwnedBy:     provider.Code,
				Provider:    provider.Name,
				Description: m.Description,
			})
		}
	}

	return response, nil
}

func (s *APIService) getActiveModel(modelCode string) (*model.Model, error) {
	var m model.Model
	if err := repository.DB.Preload("Provider").
		Where("code = ? AND status = ?", modelCode, 1).
		First(&m).Error; err != nil {
		return nil, errors.New("模型不存在或已禁用")
	}

	if m.Provider.Status != 1 {
		return nil, errors.New("模型服务商已禁用")
	}

	return &m, nil
}

func (s *APIService) checkUserPoints(userID uint) (float64, error) {
	var user model.User
	if err := repository.DB.Select("points").First(&user, userID).Error; err != nil {
		return 0, errors.New("用户不存在")
	}
	return user.Points, nil
}

func (s *APIService) estimateInputTokens(messages []model.ChatMessage) int {
	totalChars := 0
	for _, msg := range messages {
		content, ok := msg.Content.(string)
		if ok {
			totalChars += len(content)
		}
	}
	return totalChars / 4
}

func (s *APIService) estimateCost(model *model.Model, inputTokens, outputTokens int) float64 {
	return s.calculatePoints(model, inputTokens, outputTokens) * 2
}

func (s *APIService) calculatePoints(model *model.Model, inputTokens, outputTokens int) float64 {
	inputPoints := (float64(inputTokens) / 1000.0) * model.PointsPer1KInput
	outputPoints := (float64(outputTokens) / 1000.0) * model.PointsPer1KOutput
	return inputPoints + outputPoints
}

func (s *APIService) callModelAPI(modelInfo *model.Model, req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, int, int, error) {
	provider := modelInfo.Provider

	switch provider.Code {
	case "openai":
		return s.callOpenAI(provider, modelInfo, req)
	case "anthropic":
		return s.callAnthropic(provider, modelInfo, req)
	case "zhipu":
		return s.callZhipu(provider, modelInfo, req)
	case "qwen":
		return s.callQwen(provider, modelInfo, req)
	case "hunyuan":
		return s.callHunyuan(provider, modelInfo, req)
	case "doubao":
		return s.callDoubao(provider, modelInfo, req)
	default:
		return s.callMockModel(modelInfo, req)
	}
}

func (s *APIService) callMockModel(modelInfo *model.Model, req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, int, int, error) {
	inputTokens := s.estimateInputTokens(req.Messages)
	outputTokens := 50

	response := &model.ChatCompletionResponse{
		ID:      "mock-" + uuid.New().String(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   modelInfo.Code,
		Choices: []struct {
			Index        int         `json:"index"`
			Message      model.ChatMessage `json:"message"`
			FinishReason string      `json:"finish_reason"`
		}{
			{
				Index: 0,
				Message: model.ChatMessage{
					Role:    "assistant",
					Content: fmt.Sprintf("这是来自模型 %s 的模拟响应。您发送了 %d 条消息，共约 %d 个tokens。", modelInfo.Name, len(req.Messages), inputTokens),
				},
				FinishReason: "stop",
			},
		},
	}

	response.Usage.PromptTokens = inputTokens
	response.Usage.CompletionTokens = outputTokens
	response.Usage.TotalTokens = inputTokens + outputTokens

	return response, inputTokens, outputTokens, nil
}

func (s *APIService) callOpenAI(provider model.Provider, modelInfo *model.Model, req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, int, int, error) {
	if provider.APIKey == "" {
		return s.callMockModel(modelInfo, req)
	}

	body, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", provider.APIEndpoint+"/v1/chat/completions", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+provider.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, 0, 0, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, 0, 0, errors.New(string(respBody))
	}

	var response model.ChatCompletionResponse
	json.Unmarshal(respBody, &response)

	return &response, response.Usage.PromptTokens, response.Usage.CompletionTokens, nil
}

func (s *APIService) callAnthropic(provider model.Provider, modelInfo *model.Model, req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, int, int, error) {
	return s.callMockModel(modelInfo, req)
}

func (s *APIService) callZhipu(provider model.Provider, modelInfo *model.Model, req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, int, int, error) {
	return s.callMockModel(modelInfo, req)
}

func (s *APIService) callQwen(provider model.Provider, modelInfo *model.Model, req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, int, int, error) {
	return s.callMockModel(modelInfo, req)
}

func (s *APIService) callHunyuan(provider model.Provider, modelInfo *model.Model, req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, int, int, error) {
	return s.callMockModel(modelInfo, req)
}

func (s *APIService) callDoubao(provider model.Provider, modelInfo *model.Model, req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, int, int, error) {
	return s.callMockModel(modelInfo, req)
}
