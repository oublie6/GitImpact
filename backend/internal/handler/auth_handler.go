package handler

import (
	"errors"

	"gitimpact/backend/internal/service"
	"gitimpact/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

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
func (h *AuthHandler) Logout(c *gin.Context) { response.OK(c, gin.H{"message": "logout success"}) }
func (h *AuthHandler) Me(c *gin.Context) {
	response.OK(c, gin.H{"username": c.GetString("username"), "user_id": c.GetUint("user_id"), "role": c.GetString("role")})
}
