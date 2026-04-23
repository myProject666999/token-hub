package handler

import (
	"token-hub/internal/model"
	"token-hub/internal/service"
	"token-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	apiService    *service.APIService
	apiKeyService *service.APIKeyService
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{
		apiService:    service.NewAPIService(),
		apiKeyService: service.NewAPIKeyService(),
	}
}

func (h *APIHandler) ListModels(c *gin.Context) {
	userID := c.GetUint("user_id")

	models, err := h.apiService.GetAvailableModels(userID)
	if err != nil {
		response.InternalServerError(c, "获取模型列表失败")
		return
	}

	response.Success(c, models)
}

func (h *APIHandler) ChatCompletion(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req model.ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Model == "" {
		response.BadRequest(c, "模型名称不能为空")
		return
	}

	if len(req.Messages) == 0 {
		response.BadRequest(c, "消息列表不能为空")
		return
	}

	resp, err := h.apiService.ChatCompletion(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(200, resp)
}

func (h *APIHandler) Completion(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req model.CompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Model == "" {
		response.BadRequest(c, "模型名称不能为空")
		return
	}

	if req.Prompt == "" {
		response.BadRequest(c, "提示内容不能为空")
		return
	}

	resp, err := h.apiService.Completion(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(200, resp)
}

func (h *APIHandler) CreateAPIKey(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req model.APIKeyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	key, err := h.apiKeyService.CreateAPIKey(userID, req.Name)
	if err != nil {
		response.InternalServerError(c, "创建API密钥失败")
		return
	}

	response.Success(c, key)
}

func (h *APIHandler) GetAPIKeys(c *gin.Context) {
	userID := c.GetUint("user_id")

	keys, err := h.apiKeyService.GetAPIKeys(userID)
	if err != nil {
		response.InternalServerError(c, "获取API密钥失败")
		return
	}

	response.Success(c, keys)
}

func (h *APIHandler) DeleteAPIKey(c *gin.Context) {
	userID := c.GetUint("user_id")
	keyIDStr := c.Param("id")
	keyID, err := strconvUint(keyIDStr)
	if err != nil {
		response.BadRequest(c, "无效的密钥ID")
		return
	}

	if err := h.apiKeyService.DeleteAPIKey(userID, keyID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func (h *APIHandler) UpdateAPIKeyStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	keyIDStr := c.Param("id")
	keyID, err := strconvUint(keyIDStr)
	if err != nil {
		response.BadRequest(c, "无效的密钥ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if req.Status != 0 && req.Status != 1 {
		response.BadRequest(c, "无效的状态值")
		return
	}

	if err := h.apiKeyService.UpdateAPIKeyStatus(userID, keyID, req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "状态更新成功", nil)
}

func strconvUint(s string) (uint, error) {
	var result uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		result = result*10 + uint(c-'0')
	}
	return result, nil
}
