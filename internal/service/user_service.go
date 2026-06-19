package service

import (
	"errors"
	"mime/multipart"
	"strings"

	"starter-api-golang/internal/config"
	"starter-api-golang/internal/domain/entity"
	domainRepo "starter-api-golang/internal/domain/repository"
	"starter-api-golang/internal/domain/usecase"
	"starter-api-golang/pkg/hash"
	"starter-api-golang/pkg/upload"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserEmailTaken    = errors.New("email already taken")
	ErrUserUsernameTaken = errors.New("username already taken")
	ErrUserNotFoundSvc   = errors.New("user not found")
)

// applyPhotoURL replaces the stored filename with a full URL, in-memory only.
// It is idempotent: already-absolute URLs (starting with "http") are left unchanged.
func applyPhotoURL(storageURL string, user *entity.User) {
	if user == nil || user.Photo == nil || *user.Photo == "" {
		return
	}
	if strings.HasPrefix(*user.Photo, "http") {
		return
	}
	url := upload.BuildPhotoURL(storageURL, *user.Photo)
	user.Photo = &url
}

type userService struct {
	userRepo domainRepo.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo domainRepo.UserRepository, cfg *config.Config) usecase.UserUsecase {
	return &userService{userRepo: userRepo, cfg: cfg}
}

func (s *userService) FindAll(filter domainRepo.UserFilter) ([]entity.User, int64, error) {
	users, total, err := s.userRepo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}
	for i := range users {
		applyPhotoURL(s.cfg.Storage.URL, &users[i])
	}
	return users, total, nil
}

func (s *userService) FindByID(id uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFoundSvc
	}
	if err != nil {
		return nil, err
	}
	applyPhotoURL(s.cfg.Storage.URL, user)
	return user, nil
}

func (s *userService) Create(req usecase.CreateUserRequest) (*entity.User, error) {
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, ErrUserEmailTaken
	}
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, ErrUserUsernameTaken
	}

	hashed, err := hash.Make(req.Password)
	if err != nil {
		return nil, err
	}

	status := true
	if req.IsActive != nil {
		status = *req.IsActive
	}

	user := &entity.User{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: &hashed,
		RoleID:   req.RoleID,
		Status:   status,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	user, err = s.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, err
	}
	applyPhotoURL(s.cfg.Storage.URL, user)
	return user, nil
}

func (s *userService) Update(id uuid.UUID, req usecase.UpdateUserRequest) (*entity.User, error) {
	user, err := s.userRepo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFoundSvc
	}
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.RoleID != nil {
		user.RoleID = req.RoleID
	}
	if req.IsActive != nil {
		user.Status = *req.IsActive
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	user, err = s.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, err
	}
	applyPhotoURL(s.cfg.Storage.URL, user)
	return user, nil
}

func (s *userService) Delete(id uuid.UUID) error {
	if _, err := s.userRepo.FindByID(id); errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFoundSvc
	}
	return s.userRepo.Delete(id)
}

func (s *userService) UpdateProfile(userID uuid.UUID, req usecase.UpdateProfileRequest) (*entity.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFoundSvc
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Username != "" && req.Username != user.Username {
		if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
			return nil, ErrUserUsernameTaken
		}
		user.Username = req.Username
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	applyPhotoURL(s.cfg.Storage.URL, user)
	return user, nil
}

func (s *userService) UpdatePhoto(userID uuid.UUID, file *multipart.FileHeader) (*entity.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFoundSvc
	}

	if user.Photo != nil && *user.Photo != "" {
		_ = upload.DeletePhoto(s.cfg.Storage.Path, *user.Photo)
	}

	filename, err := upload.SavePhoto(file, s.cfg.Storage.Path)
	if err != nil {
		return nil, err
	}

	user.Photo = &filename
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	applyPhotoURL(s.cfg.Storage.URL, user)
	return user, nil
}
