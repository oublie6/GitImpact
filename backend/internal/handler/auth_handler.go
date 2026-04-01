// auth_handler.go 提供注册、登录、退出和当前用户查询接口。
package handler

import (
	"errors"

	"gitimpact/backend/internal/service"
	"gitimpact/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 是认证相关 HTTP 入口。
type AuthHandler struct{ svc *service.AuthService }

// NewAuthHandler 创建认证处理器。
func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

// Register 处理数据库用户注册请求。
// 请求体仅包含用户名、密码和邮箱，成功后返回基础确认信息。
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct{ Username, Password, Email string }
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	if err := h.svc.Register(req.Username, req.Password, req.Email); err != nil {
		if errors.Is(err, service.ErrUserExists) {
			response.Err(c, 409, err.Error())
			return
		}
		response.Err(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"username": req.Username})
}

// Login 校验用户身份并签发 JWT。
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct{ Username, Password string }
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	u, err := h.svc.ValidateLogin(req.Username, req.Password)
	if err != nil {
		response.Err(c, 401, err.Error())
		return
	}
	token, err := h.svc.GenerateToken(u)
	if err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, gin.H{"token": token, "user": u})
}

// Logout 当前仅作为前端状态清理的确认接口，后端未实现令牌黑名单。
func (h *AuthHandler) Logout(c *gin.Context) { response.OK(c, gin.H{"message": "logout success"}) }

// Me 返回 JWT 中解析出的当前用户上下文。
// 对配置用户而言，user_id 可能为 0，这是当前实现的真实行为。
func (h *AuthHandler) Me(c *gin.Context) {
	response.OK(c, gin.H{"username": c.GetString("username"), "user_id": c.GetUint("user_id"), "role": c.GetString("role")})
}
