package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/dto"
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type AuthHandler struct {
	App *app.App
}

func NewAuthHandler(app *app.App) *AuthHandler {
	return &AuthHandler{
		App: app,
	}
}

type RegisterRequest struct {
	UserName string         `json:"username" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Role     model.UserRole `json:"user_role" binding:"required"`
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	if err := h.App.UserService.CreateUser(req.UserName, req.Password); err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			dto.Conflict(ctx, "USER_CONFLICTS", fmt.Sprintf("User named %s already exists", req.UserName))
			return
		}
		dto.InternalServerError(ctx, "Failed to create user")
		return
	}

	dto.Success(ctx, 201, fmt.Sprintf("Created user named %s successfully", req.UserName))
}

type LoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	tokenStr, err := h.App.AuthService.Login(req.UserName, req.Password)
	if err != nil {
		dto.InternalServerError(ctx, "Failed to login")
	}
	// change the parameter secure to true when deploy
	ctx.SetCookie("jwt", tokenStr, 3600, "/", "", false, true)

	dto.Success(ctx, http.StatusOK, gin.H{
		"status": "Login successfully",
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "/", "", false, true)
	dto.Success(ctx, http.StatusOK, "Logged out successfully")
}
