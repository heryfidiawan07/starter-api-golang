package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"starter-api-golang/internal/config"
	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"
	"starter-api-golang/internal/domain/usecase"
	"starter-api-golang/pkg/email"
	"starter-api-golang/pkg/hash"
	"starter-api-golang/pkg/jwt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already taken")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrEmailNotVerified   = errors.New("email not verified, please check your inbox")
	ErrAccountInactive    = errors.New("account is inactive")
	ErrWrongPassword      = errors.New("current password is incorrect")
	ErrSocialNoEmail      = errors.New("could not retrieve email from social provider")
)

type authService struct {
	userRepo         domainRepo.UserRepository
	refreshTokenRepo domainRepo.RefreshTokenRepository
	resetTokenRepo   domainRepo.PasswordResetTokenRepository
	socialRepo       domainRepo.SocialAccountRepository
	mailer           *email.Mailer
	cfg              *config.Config
}

func NewAuthService(
	userRepo domainRepo.UserRepository,
	refreshTokenRepo domainRepo.RefreshTokenRepository,
	resetTokenRepo domainRepo.PasswordResetTokenRepository,
	socialRepo domainRepo.SocialAccountRepository,
	mailer *email.Mailer,
	cfg *config.Config,
) usecase.AuthUsecase {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		resetTokenRepo:   resetTokenRepo,
		socialRepo:       socialRepo,
		mailer:           mailer,
		cfg:              cfg,
	}
}

func (s *authService) Register(req usecase.RegisterRequest) (*entity.User, error) {
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, ErrEmailTaken
	}
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, ErrUsernameTaken
	}

	hashed, err := hash.Make(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: &hashed,
		Status:   true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	if s.cfg.EmailVerificationRequired {
		go s.sendVerificationEmail(user)
	} else {
		now := time.Now()
		user.VerifiedAt = &now
		_ = s.userRepo.Update(user)
	}

	applyPhotoURL(s.cfg.Storage.URL, user)
	return user, nil
}

func (s *authService) Login(req usecase.LoginRequest) (*usecase.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.Status {
		return nil, ErrAccountInactive
	}

	if user.Password == nil || !hash.Check(req.Password, *user.Password) {
		return nil, ErrInvalidCredentials
	}

	if s.cfg.EmailVerificationRequired && user.VerifiedAt == nil {
		return nil, ErrEmailNotVerified
	}

	return s.generateTokenPair(user)
}

func (s *authService) Logout(userID uuid.UUID, refreshToken string) error {
	return s.refreshTokenRepo.RevokeAllByUserID(userID)
}

func (s *authService) RefreshToken(token string) (*usecase.AuthResponse, error) {
	rt, err := s.refreshTokenRepo.FindByToken(token)
	if err != nil || !rt.IsValid() {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(rt.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := s.refreshTokenRepo.Revoke(token); err != nil {
		return nil, err
	}

	return s.generateTokenPair(user)
}

func (s *authService) RevokeToken(refreshToken string) error {
	rt, err := s.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		return ErrInvalidToken
	}
	if !rt.IsValid() {
		return ErrInvalidToken
	}
	return s.refreshTokenRepo.Revoke(refreshToken)
}

func (s *authService) ForgotPassword(req usecase.ForgotPasswordRequest) error {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil // don't reveal if email exists
	}

	_ = s.resetTokenRepo.DeleteByUserID(user.ID)

	token, err := generateSecureToken(32)
	if err != nil {
		return err
	}

	expiry := time.Now().Add(time.Duration(s.cfg.ResetTokenExpire) * time.Minute)
	prt := &entity.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiry,
	}
	if err := s.resetTokenRepo.Create(prt); err != nil {
		return err
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.App.URL, token)
	go func() {
		_ = s.mailer.SendPasswordReset(user.Email, user.Name, resetURL)
	}()

	return nil
}

func (s *authService) ResetPassword(req usecase.ResetPasswordRequest) error {
	prt, err := s.resetTokenRepo.FindByToken(req.Token)
	if err != nil || !prt.IsValid() {
		return ErrInvalidToken
	}

	hashed, err := hash.Make(req.Password)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(prt.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	user.Password = &hashed
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	_ = s.resetTokenRepo.MarkUsed(req.Token)
	_ = s.refreshTokenRepo.RevokeAllByUserID(user.ID)
	return nil
}

func (s *authService) VerifyEmail(token string) error {
	prt, err := s.resetTokenRepo.FindByToken(token)
	if err != nil || !prt.IsValid() {
		return ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(prt.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	now := time.Now()
	user.VerifiedAt = &now
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	_ = s.resetTokenRepo.MarkUsed(token)
	return nil
}

func (s *authService) GoogleAuth(req usecase.SocialAuthRequest) (*usecase.AuthResponse, error) {
	// The frontend sends an access_token from Google Identity Services SDK — use it directly.
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + req.AccessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(body, &info); err != nil || info.Email == "" {
		return nil, ErrSocialNoEmail
	}

	return s.handleSocialLogin(entity.ProviderGoogle, info.ID, info.Email, info.Name)
}

func (s *authService) FacebookAuth(req usecase.SocialAuthRequest) (*usecase.AuthResponse, error) {
	resp, err := http.Get(fmt.Sprintf(
		"https://graph.facebook.com/me?fields=id,name,email&access_token=%s",
		req.AccessToken,
	))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &info); err != nil || info.Email == "" {
		return nil, ErrSocialNoEmail
	}

	return s.handleSocialLogin(entity.ProviderFacebook, info.ID, info.Email, info.Name)
}

func (s *authService) GetMe(userID uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	applyPhotoURL(s.cfg.Storage.URL, user)
	return user, nil
}

func (s *authService) ChangePassword(userID uuid.UUID, req usecase.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if user.Password == nil || !hash.Check(req.OldPassword, *user.Password) {
		return ErrWrongPassword
	}

	hashed, err := hash.Make(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = &hashed
	return s.userRepo.Update(user)
}

func (s *authService) generateTokenPair(user *entity.User) (*usecase.AuthResponse, error) {
	applyPhotoURL(s.cfg.Storage.URL, user)
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.IsRoot, s.cfg.JWT.Secret, s.cfg.JWT.AccessExpire)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := jwt.GenerateRefreshToken(user.ID, s.cfg.JWT.Secret, s.cfg.JWT.RefreshExpire)
	if err != nil {
		return nil, err
	}

	rt := &entity.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.JWT.RefreshExpire) * time.Minute),
	}
	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return nil, err
	}

	return &usecase.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    s.cfg.JWT.AccessExpire * 60,
		User:         user,
	}, nil
}

func (s *authService) handleSocialLogin(provider entity.SocialProvider, providerID, providerEmail, name string) (*usecase.AuthResponse, error) {
	account, err := s.socialRepo.FindByProviderAndID(provider, providerID)

	var user *entity.User

	if err != nil {
		// no existing social account — find or create user
		user, err = s.userRepo.FindByEmail(providerEmail)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			username := generateUsername(providerEmail)
			now := time.Now()
			user = &entity.User{
				Username:   username,
				Name:       name,
				Email:      providerEmail,
				Status:     true,
				VerifiedAt: &now,
			}
			if err := s.userRepo.Create(user); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		sa := &entity.SocialAccount{
			UserID:        user.ID,
			Provider:      provider,
			ProviderID:    providerID,
			ProviderEmail: providerEmail,
		}
		_ = s.socialRepo.Create(sa)
	} else {
		user, err = s.userRepo.FindByID(account.UserID)
		if err != nil {
			return nil, ErrUserNotFound
		}
	}

	if !user.Status {
		return nil, ErrAccountInactive
	}

	return s.generateTokenPair(user)
}

func (s *authService) sendVerificationEmail(user *entity.User) {
	token, err := generateSecureToken(32)
	if err != nil {
		return
	}

	expiry := time.Now().Add(24 * time.Hour)
	prt := &entity.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiry,
	}
	if err := s.resetTokenRepo.Create(prt); err != nil {
		return
	}

	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", s.cfg.App.URL, token)
	_ = s.mailer.SendVerification(user.Email, user.Name, verifyURL)
}

func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func generateUsername(email string) string {
	for i, c := range email {
		if c == '@' {
			return email[:i] + "_" + fmt.Sprintf("%d", time.Now().Unix())[:4]
		}
	}
	return "user_" + fmt.Sprintf("%d", time.Now().UnixNano())
}
