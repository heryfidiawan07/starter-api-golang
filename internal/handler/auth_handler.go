package handler

import (
	"errors"

	"starter-api-golang/internal/domain/usecase"
	"starter-api-golang/internal/middleware"
	"starter-api-golang/internal/service"
	"starter-api-golang/pkg/response"
	"starter-api-golang/pkg/validator"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req usecase.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	user, err := h.authUsecase.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailTaken):
			response.BadRequest(c, err.Error(), nil)
		case errors.Is(err, service.ErrUsernameTaken):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalServerError(c, "registration failed")
		}
		return
	}

	response.Created(c, "registration successful, please verify your email", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	result, err := h.authUsecase.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Unauthorized(c, err.Error())
		case errors.Is(err, service.ErrEmailNotVerified):
			response.Unauthorized(c, err.Error())
		case errors.Is(err, service.ErrAccountInactive):
			response.Forbidden(c, err.Error())
		default:
			response.InternalServerError(c, "login failed")
		}
		return
	}

	response.OK(c, "login successful", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req usecase.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.authUsecase.Logout(userID, req.RefreshToken); err != nil {
		response.InternalServerError(c, "logout failed")
		return
	}

	response.OK(c, "logout successful", nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req usecase.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	result, err := h.authUsecase.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid or expired refresh token")
		return
	}

	response.OK(c, "token refreshed", result)
}

func (h *AuthHandler) RevokeToken(c *gin.Context) {
	var req usecase.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if err := h.authUsecase.RevokeToken(req.RefreshToken); err != nil {
		response.BadRequest(c, "invalid token", nil)
		return
	}

	response.OK(c, "token revoked", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req usecase.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	_ = h.authUsecase.ForgotPassword(req)
	response.OK(c, "if the email exists, a reset link has been sent", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req usecase.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	if err := h.authUsecase.ResetPassword(req); err != nil {
		response.BadRequest(c, "invalid or expired reset token", nil)
		return
	}

	response.OK(c, "password reset successful", nil)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.BadRequest(c, "token is required", nil)
		return
	}

	if err := h.authUsecase.VerifyEmail(token); err != nil {
		response.BadRequest(c, "invalid or expired verification token", nil)
		return
	}

	response.OK(c, "email verified successfully", nil)
}

func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	var req usecase.SocialAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	result, err := h.authUsecase.GoogleAuth(req)
	if err != nil {
		response.BadRequest(c, "google authentication failed: "+err.Error(), nil)
		return
	}

	response.OK(c, "google login successful", result)
}

func (h *AuthHandler) FacebookAuth(c *gin.Context) {
	var req usecase.SocialAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	result, err := h.authUsecase.FacebookAuth(req)
	if err != nil {
		response.BadRequest(c, "facebook authentication failed: "+err.Error(), nil)
		return
	}

	response.OK(c, "facebook login successful", result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.authUsecase.GetMe(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.OK(c, "profile retrieved", user)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req usecase.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.authUsecase.ChangePassword(userID, req); err != nil {
		switch {
		case errors.Is(err, service.ErrWrongPassword):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalServerError(c, "change password failed")
		}
		return
	}

	response.OK(c, "password changed successfully", nil)
}
