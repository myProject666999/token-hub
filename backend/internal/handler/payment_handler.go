package handler

import (
	"strconv"
	"token-hub/internal/model"
	"token-hub/internal/service"
	"token-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: service.NewPaymentService(),
	}
}

func (h *PaymentHandler) GetPaymentMethods(c *gin.Context) {
	includeDisabled := c.Query("include_disabled") == "true"

	methods, err := h.paymentService.GetPaymentMethods(includeDisabled)
	if err != nil {
		response.InternalServerError(c, "获取支付方式失败")
		return
	}

	response.Success(c, methods)
}

func (h *PaymentHandler) CreatePaymentMethod(c *gin.Context) {
	var method model.PaymentMethod
	if err := c.ShouldBindJSON(&method); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.paymentService.CreatePaymentMethod(&method); err != nil {
		response.InternalServerError(c, "创建支付方式失败")
		return
	}

	response.Success(c, method)
}

func (h *PaymentHandler) UpdatePaymentMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var method model.PaymentMethod
	if err := c.ShouldBindJSON(&method); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.paymentService.UpdatePaymentMethod(uint(id), &method); err != nil {
		response.InternalServerError(c, "更新支付方式失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

func (h *PaymentHandler) DeletePaymentMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.paymentService.DeletePaymentMethod(uint(id)); err != nil {
		response.InternalServerError(c, "删除支付方式失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func (h *PaymentHandler) CreateRechargeOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req model.RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Amount <= 0 {
		response.BadRequest(c, "充值金额必须大于0")
		return
	}

	record, err := h.paymentService.CreateRechargeOrder(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, record)
}

func (h *PaymentHandler) SimulatePayment(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		response.BadRequest(c, "订单号不能为空")
		return
	}

	if err := h.paymentService.SimulatePayment(orderNo); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "模拟支付成功", nil)
}

func (h *PaymentHandler) GetUserRechargeRecords(c *gin.Context) {
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

	records, total, err := h.paymentService.GetRechargeRecordList(userID, page, pageSize, status)
	if err != nil {
		response.InternalServerError(c, "获取充值记录失败")
		return
	}

	response.SuccessPage(c, records, total, page, pageSize)
}

func (h *PaymentHandler) GetAllRechargeRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	records, total, err := h.paymentService.GetAllRechargeRecords(page, pageSize, status)
	if err != nil {
		response.InternalServerError(c, "获取充值记录失败")
		return
	}

	response.SuccessPage(c, records, total, page, pageSize)
}
