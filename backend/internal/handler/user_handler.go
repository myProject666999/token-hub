package handler

import (
	"strconv"
	"token-hub/internal/model"
	"token-hub/internal/service"
	"token-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService   *service.UserService
	pointsService *service.PointsService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService:   service.NewUserService(),
		pointsService: service.NewPointsService(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	token, user, err := h.userService.Login(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userService.GetByID(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Update(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetPoints(c *gin.Context) {
	userID := c.GetUint("user_id")

	points, err := h.userService.GetUserPoints(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, points)
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := h.userService.GetList(page, pageSize, keyword)
	if err != nil {
		response.InternalServerError(c, "获取用户列表失败")
		return
	}

	response.SuccessPage(c, users, total, page, pageSize)
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
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

	if err := h.userService.UpdateStatus(uint(userID), req.Status); err != nil {
		response.InternalServerError(c, "更新用户状态失败")
		return
	}

	response.SuccessWithMessage(c, "状态更新成功", nil)
}

func (h *UserHandler) GetUserStatistics(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := h.pointsService.GetUserStatistics(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, stats)
}
