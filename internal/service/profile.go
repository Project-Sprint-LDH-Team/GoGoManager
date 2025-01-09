package service

import (
	"context"
	"errors"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/repository"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID uint) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uint, req *models.UpdateProfileRequest) (*models.User, error)
}

type profileService struct {
	userRepo repository.UserRepository
}

func NewProfileService(userRepo repository.UserRepository) ProfileService {
	return &profileService{
		userRepo: userRepo,
	}
}

func (s *profileService) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID uint, req *models.UpdateProfileRequest) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if email is changed and already exists
	if req.Email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.New("email is already used by another user")
		}
	}

	// Update user fields
	user.Email = req.Email
	user.Name = req.Name
	user.UserImageUri = req.UserImageUri
	user.CompanyName = req.CompanyName
	user.CompanyImageUri = req.CompanyImageUri

	// Save updates
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
