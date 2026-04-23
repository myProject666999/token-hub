package handler

import (
	"strconv"
	"token-hub/internal/service"
	"token-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

type PointsHandler struct {
	pointsService *service.PointsService
}

func NewPointsHandler() *PointsHandler {
	return &PointsHandler{
		pointsService: service.NewPointsService(),
	}
}

func (h *PointsHandler) GetPointsRate(c *gin.Context) {
	rate, err := h.pointsService.GetPointsRate()
	if err != nil {
		response.InternalServerError(c, "获取积分汇率失败")
		return
	}

	response.Success(c, gin.H{
		"rate": rate,
	})
}

func (h *PointsHandler) SetPointsRate(c *gin.Context) {
	var req struct {
		Rate float64 `json:"rate" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.pointsService.SetPointsRate(req.Rate); err != nil {
		response.InternalServerError(c, "设置积分汇率失败")
		return
	}

	response.SuccessWithMessage(c, "设置成功", nil)
}

func (h *PointsHandler) GetPointsRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	recordType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	records, total, err := h.pointsService.GetPointsRecordList(userID, page, pageSize, recordType)
	if err != nil {
		response.InternalServerError(c, "获取积分记录失败")
		return
	}

	response.SuccessPage(c, records, total, page, pageSize)
}

func (h *PointsHandler) GetAllStatistics(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	stats, total, err := h.pointsService.GetAllStatistics(page, pageSize)
	if err != nil {
		response.InternalServerError(c, "获取统计信息失败")
		return
	}

	response.SuccessPage(c, stats, total, page, pageSize)
}
