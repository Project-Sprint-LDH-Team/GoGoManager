package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Email           string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password        string    `gorm:"size:255;not null" json:"-"`
	Name            string    `gorm:"size:52" json:"name"`
	UserImageUri    string    `gorm:"size:255" json:"user_image_uri"`
	CompanyName     string    `gorm:"size:52" json:"company_name"`
	CompanyImageUri string    `gorm:"size:255" json:"company_image_uri"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relations
	Departments []Department `gorm:"foreignKey:UserID" json:"-"`
	Files       []File       `gorm:"foreignKey:UserID" json:"-"`
	Employees   []Employee   `gorm:"foreignKey:UserID" json:"-"`
}

// Request structs
type AuthRequest struct {
	Email    string `json:"email" validate:"required,email,min=3,max=255"`
	Password string `json:"password" validate:"required,min=8,max=32"`
	Action   string `json:"action" validate:"required,oneof=create login"`
}

// Response structs
type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// UpdateProfileRequest untuk PATCH /v1/user
type UpdateProfileRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Name            string `json:"name" validate:"required,min=4,max=52"`
	UserImageUri    string `json:"userImageUri" validate:"omitempty,uri"`
	CompanyName     string `json:"companyName" validate:"required,min=4,max=52"`
	CompanyImageUri string `json:"companyImageUri" validate:"omitempty,uri"`
}

// ProfileResponse untuk GET & PATCH /v1/user response
type ProfileResponse struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	UserImageUri    string `json:"userImageUri"`
	CompanyName     string `json:"companyName"`
	CompanyImageUri string `json:"companyImageUri"`
}

// Method untuk convert User ke ProfileResponse
func (u *User) ToProfileResponse() *ProfileResponse {
	return &ProfileResponse{
		Email:           u.Email,
		Name:            u.Name,
		UserImageUri:    u.UserImageUri,
		CompanyName:     u.CompanyName,
		CompanyImageUri: u.CompanyImageUri,
	}
}

// BeforeCreate hook untuk hash password
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
