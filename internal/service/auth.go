package service

import (
	"context"
	"errors"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/repository"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/pkg/jwt"
)

type AuthService interface {
	Authenticate(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtMaker jwt.Maker
}

func NewAuthService(userRepo repository.UserRepository, jwtMaker jwt.Maker) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtMaker: jwtMaker,
	}
}

func (s *authService) Authenticate(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error) {
	switch req.Action {
	case "create":
		return s.register(ctx, req)
	case "login":
		return s.login(ctx, req)
	default:
		return nil, errors.New("invalid action")
	}
}

func (s *authService) register(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error) {
	// Check if email exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create new user
	user := &models.User{
		Email:    req.Email,
		Password: req.Password, // Password will be hashed by GORM hook BeforeCreate
	}

	// Save user (GORM will automatically hash the password via BeforeCreate hook)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.jwtMaker.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Email: user.Email,
		Token: token,
	}, nil
}

func (s *authService) login(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Compare password
	if err := user.ComparePassword(req.Password); err != nil {
		return nil, errors.New("invalid password")
	}

	// Generate token
	token, err := s.jwtMaker.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Email: user.Email,
		Token: token,
	}, nil
}
