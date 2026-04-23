package middleware

import (
	"strings"
	"token-hub/internal/model"
	"token-hub/internal/repository"
	"token-hub/pkg/jwt"
	"token-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		var user model.User
		if err := repository.DB.First(&user, claims.UserID).Error; err != nil {
			response.Unauthorized(c, "用户不存在")
			c.Abort()
			return
		}

		if user.Status != 1 {
			response.Forbidden(c, "用户已被禁用")
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)

		c.Next()
	}
}

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		if role != "admin" {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供API密钥")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "API密钥格式错误")
			c.Abort()
			return
		}

		apiKey := parts[1]

		var key model.APIKey
		if err := repository.DB.Where("key = ?", apiKey).First(&key).Error; err != nil {
			response.Unauthorized(c, "API密钥无效")
			c.Abort()
			return
		}

		if key.Status != 1 {
			response.Forbidden(c, "API密钥已被禁用")
			c.Abort()
			return
		}

		var user model.User
		if err := repository.DB.First(&user, key.UserID).Error; err != nil {
			response.Unauthorized(c, "用户不存在")
			c.Abort()
			return
		}

		if user.Status != 1 {
			response.Forbidden(c, "用户已被禁用")
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)
		c.Set("api_key_id", key.ID)

		c.Next()
	}
}
