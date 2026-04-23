package handler

import (
	"strconv"
	"token-hub/internal/model"
	"token-hub/internal/service"
	"token-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	providerService *service.ProviderService
}

func NewProviderHandler() *ProviderHandler {
	return &ProviderHandler{
		providerService: service.NewProviderService(),
	}
}

func (h *ProviderHandler) GetAllProviders(c *gin.Context) {
	includeDisabled := c.Query("include_disabled") == "true"

	providers, err := h.providerService.GetAllProviders(includeDisabled)
	if err != nil {
		response.InternalServerError(c, "获取服务商列表失败")
		return
	}

	response.Success(c, providers)
}

func (h *ProviderHandler) GetProviderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	providers, total, err := h.providerService.GetProviderList(page, pageSize, keyword)
	if err != nil {
		response.InternalServerError(c, "获取服务商列表失败")
		return
	}

	response.SuccessPage(c, providers, total, page, pageSize)
}

func (h *ProviderHandler) GetProviderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的服务商ID")
		return
	}

	provider, err := h.providerService.GetProviderByID(uint(id))
	if err != nil {
		response.NotFound(c, "服务商不存在")
		return
	}

	response.Success(c, provider)
}

func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var req model.ProviderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	provider, err := h.providerService.CreateProvider(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, provider)
}

func (h *ProviderHandler) UpdateProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的服务商ID")
		return
	}

	var req model.ProviderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	provider, err := h.providerService.UpdateProvider(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, provider)
}

func (h *ProviderHandler) DeleteProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的服务商ID")
		return
	}

	if err := h.providerService.DeleteProvider(uint(id)); err != nil {
		response.InternalServerError(c, "删除服务商失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func (h *ProviderHandler) GetModelList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	providerIDStr := c.Query("provider_id")
	var providerID uint
	if providerIDStr != "" {
		id, _ := strconv.ParseUint(providerIDStr, 10, 64)
		providerID = uint(id)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	models, total, err := h.providerService.GetModelList(page, pageSize, providerID, keyword)
	if err != nil {
		response.InternalServerError(c, "获取模型列表失败")
		return
	}

	response.SuccessPage(c, models, total, page, pageSize)
}

func (h *ProviderHandler) GetModelByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的模型ID")
		return
	}

	model, err := h.providerService.GetModelByID(uint(id))
	if err != nil {
		response.NotFound(c, "模型不存在")
		return
	}

	response.Success(c, model)
}

func (h *ProviderHandler) CreateModel(c *gin.Context) {
	var req model.ModelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	model, err := h.providerService.CreateModel(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, model)
}

func (h *ProviderHandler) UpdateModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的模型ID")
		return
	}

	var req model.ModelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	model, err := h.providerService.UpdateModel(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, model)
}

func (h *ProviderHandler) DeleteModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的模型ID")
		return
	}

	if err := h.providerService.DeleteModel(uint(id)); err != nil {
		response.InternalServerError(c, "删除模型失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
