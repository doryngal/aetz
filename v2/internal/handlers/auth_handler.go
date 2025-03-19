package handlers

import (
	"binai.net/v2/config"
	"binai.net/v2/internal/models"
	"binai.net/v2/internal/repository"
	"binai.net/v2/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(repo repository.AuthRepository, cfg *config.Config) *AuthHandler {
	service := services.NewAuthService(repo, cfg)
	return &AuthHandler{Service: service}
}

func (ac *AuthHandler) Register(c *gin.Context) {
	var registerData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user := models.User{
		Email:        registerData.Email,
		PasswordHash: registerData.Password,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ac.Service.Register(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (ac *AuthHandler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	token, role, err := ac.Service.Login(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "role": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "role": role})
}

func (ac *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.Service.ConfirmationCode(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset code sent to your email"})
}

func (ac *AuthHandler) ConfirmationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, role, err := ac.Service.ResetPassword(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "role": role})
}
