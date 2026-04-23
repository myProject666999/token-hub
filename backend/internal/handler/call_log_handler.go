package handler

import (
	"strconv"
	"token-hub/internal/service"
	"token-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

type CallLogHandler struct {
	callLogService *service.CallLogService
}

func NewCallLogHandler() *CallLogHandler {
	return &CallLogHandler{
		callLogService: service.NewCallLogService(),
	}
}

func (h *CallLogHandler) GetUserCallLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.callLogService.GetUserCallLogs(userID, page, pageSize, status)
	if err != nil {
		response.InternalServerError(c, "获取调用日志失败")
		return
	}

	response.SuccessPage(c, logs, total, page, pageSize)
}

func (h *CallLogHandler) GetAllCallLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	userIDStr := c.Query("user_id")
	var userID uint
	if userIDStr != "" {
		id, _ := strconv.ParseUint(userIDStr, 10, 64)
		userID = uint(id)
	}

	providerIDStr := c.Query("provider_id")
	var providerID uint
	if providerIDStr != "" {
		id, _ := strconv.ParseUint(providerIDStr, 10, 64)
		providerID = uint(id)
	}

	modelIDStr := c.Query("model_id")
	var modelID uint
	if modelIDStr != "" {
		id, _ := strconv.ParseUint(modelIDStr, 10, 64)
		modelID = uint(id)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.callLogService.GetAllCallLogs(page, pageSize, userID, providerID, modelID, status)
	if err != nil {
		response.InternalServerError(c, "获取调用日志失败")
		return
	}

	response.SuccessPage(c, logs, total, page, pageSize)
}

func (h *CallLogHandler) GetUserStatistics(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := h.callLogService.GetUserStatistics(userID)
	if err != nil {
		response.InternalServerError(c, "获取统计信息失败")
		return
	}

	response.Success(c, stats)
}

func (h *CallLogHandler) GetDailyStatistics(c *gin.Context) {
	userIDStr := c.Query("user_id")
	var userID uint
	if userIDStr != "" {
		id, _ := strconv.ParseUint(userIDStr, 10, 64)
		userID = uint(id)
	}

	providerIDStr := c.Query("provider_id")
	var providerID uint
	if providerIDStr != "" {
		id, _ := strconv.ParseUint(providerIDStr, 10, 64)
		providerID = uint(id)
	}

	modelIDStr := c.Query("model_id")
	var modelID uint
	if modelIDStr != "" {
		id, _ := strconv.ParseUint(modelIDStr, 10, 64)
		modelID = uint(id)
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.callLogService.GetDailyStatistics(userID, providerID, modelID, startDate, endDate)
	if err != nil {
		response.InternalServerError(c, "获取统计信息失败")
		return
	}

	response.Success(c, stats)
}

func (h *CallLogHandler) GetOverviewStatistics(c *gin.Context) {
	stats, err := h.callLogService.GetOverviewStatistics()
	if err != nil {
		response.InternalServerError(c, "获取统计信息失败")
		return
	}

	response.Success(c, stats)
}
