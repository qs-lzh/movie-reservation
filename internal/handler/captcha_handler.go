package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/dto"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type CaptchaHandler struct {
	App *app.App
}

func NewCaptchaHandler(app *app.App) *CaptchaHandler {
	return &CaptchaHandler{
		App: app,
	}
}

func (h *CaptchaHandler) GenerateCaptcha(ctx *gin.Context) {
	mBase64, tBase64, cacheKey, err := h.App.CaptchaService.Generate()
	if err != nil {
		dto.InternalServerError(ctx, "Failed to generate captcha")
		return
	}
	dto.Success(ctx, 200, gin.H{
		"image": mBase64,
		"thumb": tBase64,
		"key":   cacheKey,
	})
}

func (h *CaptchaHandler) VerifyCaptcha(ctx *gin.Context) {
	fmt.Println("start running captchaHandler.VerifyCaptcha......")
	// Define struct to hold the request body
	type CaptchaVerifyRequest struct {
		Dots []service.Dot `json:"dots" binding:"required"`
		Key  string        `json:"key" binding:"required"`
	}

	var req CaptchaVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	// Verify captcha using the service with cache key
	valid, err := h.App.CaptchaService.VerifyWithKey(req.Dots, req.Key)
	if err != nil {
		dto.InternalServerError(ctx, "Failed to verify captcha: "+err.Error())
		return
	}

	// Store validation result in cache (optional: could be used for checking later)
	// We might want to clean up the cache entry after verification regardless of success
	// Or keep the result for a short duration
	h.App.Cache.SetBool(req.Key, valid)

	dto.Success(ctx, 200, gin.H{
		"success": valid,
	})
}
