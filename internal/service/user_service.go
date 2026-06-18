package service

import (
	"errors"
	"mime/multipart"

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

type userService struct {
	userRepo domainRepo.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo domainRepo.UserRepository, cfg *config.Config) usecase.UserUsecase {
	return &userService{userRepo: userRepo, cfg: cfg}
}

func (s *userService) FindAll(filter domainRepo.UserFilter) ([]entity.User, int64, error) {
	return s.userRepo.FindAll(filter)
}

func (s *userService) FindByID(id uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFoundSvc
	}
	return user, err
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
	if req.Status != nil {
		status = *req.Status
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

	return s.userRepo.FindByID(user.ID)
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
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(user.ID)
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

	photoURL := upload.BuildPhotoURL(s.cfg.Storage.URL, filename)
	user.Photo = &photoURL
	return user, nil
}
