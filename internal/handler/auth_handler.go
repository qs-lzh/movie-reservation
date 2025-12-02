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
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	dto.Success(ctx, http.StatusOK, fmt.Sprintf("Created user named %s successfully", req.UserName))
}

type LoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	isValid, err := h.App.UserService.ValidateUser(req.UserName, req.Password)
	if !isValid {
		if err != nil {
			dto.InternalServerError(ctx, "Failed to validate user")
			return
		}
		dto.Unauthorized(ctx, "Password wrong or username doesn't not exist")
		return
	}

	标记为已登录

}

func (h *AuthHandler) Logout(ctx *gin.Context) {

	直接标记为未登录吗
}
