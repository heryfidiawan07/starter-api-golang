package usecase

import (
	"starter-api-golang/internal/domain/entity"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type SocialAuthRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	User         *entity.User `json:"user"`
}

type AuthUsecase interface {
	Register(req RegisterRequest) (*entity.User, error)
	Login(req LoginRequest) (*AuthResponse, error)
	Logout(userID uuid.UUID, refreshToken string) error
	RefreshToken(token string) (*AuthResponse, error)
	RevokeToken(refreshToken string) error
	ForgotPassword(req ForgotPasswordRequest) error
	ResetPassword(req ResetPasswordRequest) error
	VerifyEmail(token string) error
	GoogleAuth(req SocialAuthRequest) (*AuthResponse, error)
	FacebookAuth(req SocialAuthRequest) (*AuthResponse, error)
	GetMe(userID uuid.UUID) (*entity.User, error)
	ChangePassword(userID uuid.UUID, req ChangePasswordRequest) error
}
