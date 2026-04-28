package service

import (
	"Re_Shop/Backend/internal/modules/user/model"
	"Re_Shop/Backend/internal/modules/user/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrEmptyCredentials  = errors.New("username or password is empty")
	ErrUsernameExists    = errors.New("username already exists")
	ErrPhoneExists       = errors.New("phone already exists")
	ErrEmptyRegisterInfo = errors.New("username or password is empty")
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(username, password, nickname, avatar, phone string) (*model.User, error) {
	if username == "" || password == "" {
		return nil, ErrEmptyRegisterInfo
	}

	_, err := s.userRepo.FindByUsername(username)
	if err == nil {
		return nil, ErrUsernameExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if phone != "" {
		_, err = s.userRepo.FindByPhone(phone)
		if err == nil {
			return nil, ErrPhoneExists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Nickname: nickname,
		Avatar:   avatar,
		Phone:    phone,
		Role:     model.UserRoleNormal,
		Status:   model.UserStatusEnabled,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) comparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func (s *AuthService) Login(username, password string) (*model.User, error) {
	if username == "" || password == "" {
		return nil, ErrEmptyCredentials
	}

	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	if err := s.comparePassword(user.Password, password); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}
