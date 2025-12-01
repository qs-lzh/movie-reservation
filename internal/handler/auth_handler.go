package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
)

type AuthHandler struct {
	App *app.App
}

func NewAuthHandler(app *app.App) *AuthHandler {
	return &AuthHandler{
		App: app,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	ctx.Request.Body

}

func (h *AuthHandler) Login(ctx *gin.Context) {
	ctx.Request.Body

}
