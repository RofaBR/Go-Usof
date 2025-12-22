package handler

import (
	"net/http"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/internal/services"
	"github.com/gin-gonic/gin"
)

type CreateUserDTO struct {
	Login    string `json:"login" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	FullName string `json:"full_name,omitempty"`
}

type UpdateUserDTO struct {
	FullName *string `json:"full_name,omitempty"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Avatar   *string `json:"avatar,omitempty"`
}

type LoginUserDTO struct {
	Login    string `json:"login" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type AuthHandler struct {
	service *services.UserService
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req CreateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := &domain.User{
		Login:    req.Login,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     "user",
	}

	err := h.service.Register(c.Request.Context(), user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.Login(c.Request.Context(), req.Login, req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {

}
