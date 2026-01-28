package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthsarthi-dutt/blog-api/internal/service"
	"github.com/parthsarthi-dutt/blog-api/internal/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())

		return
	}

	err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		utils.Error(c, http.StatusConflict, "DUPLICATE", "Email already in use")

		return
	}

	utils.Created(c, nil, "User registered successfully")

}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())

		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "AUTH_INVALID", "Invalid credentials")

		return
	}

	utils.Success(c, gin.H{"token": token}, "Login successful")

}
