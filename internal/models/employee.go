package models

import (
	"time"
)

type Employee struct {
	ID               uint       `gorm:"primaryKey" json:"-"`
	UserID           uint       `gorm:"not null;index" json:"-"`
	DepartmentID     string     `gorm:"type:varchar(36);not null;index" json:"-"` // Pastikan tipe data sama
	IdentityNumber   string     `gorm:"uniqueIndex;size:33;not null" json:"identityNumber"`
	Name             string     `gorm:"size:33;not null" json:"name"`
	EmployeeImageUri string     `gorm:"size:255" json:"employeeImageUri"`
	Gender           string     `gorm:"size:6;not null" json:"gender"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `gorm:"index" json:"-"`

	// Relations
	User       User       `gorm:"foreignKey:UserID" json:"-"`
	Department Department `gorm:"foreignKey:DepartmentID;references:ID" json:"-"`
}

type CreateEmployeeRequest struct {
	IdentityNumber   string `json:"identityNumber" validate:"required,min=5,max=33"`
	Name             string `json:"name" validate:"required,min=4,max=33"`
	EmployeeImageUri string `json:"employeeImageUri" validate:"omitempty,uri"`
	Gender           string `json:"gender" validate:"required,oneof=male female"`
	DepartmentId     string `json:"departmentId" validate:"required"` // Perhatikan nama field ini
}

type UpdateEmployeeRequest struct {
	IdentityNumber   string `json:"identityNumber" validate:"required,min=5,max=33"`
	Name             string `json:"name" validate:"required,min=4,max=33"`
	EmployeeImageUri string `json:"employeeImageUri" validate:"omitempty,uri"`
	Gender           string `json:"gender" validate:"required,oneof=male female"`
	DepartmentId     string `json:"departmentId" validate:"required"`
}

type EmployeeResponse struct {
	IdentityNumber   string `json:"identityNumber"`
	Name             string `json:"name"`
	EmployeeImageUri string `json:"employeeImageUri"`
	Gender           string `json:"gender"`
	DepartmentId     string `json:"departmentId"`
}

func (e *Employee) ToResponse() *EmployeeResponse {
	return &EmployeeResponse{
		IdentityNumber:   e.IdentityNumber,
		Name:             e.Name,
		EmployeeImageUri: e.EmployeeImageUri,
		Gender:           e.Gender,
		DepartmentId:     e.DepartmentID,
	}
}

// Employee Filter untuk query parameters
type EmployeeFilter struct {
	Limit          int    `query:"limit"`
	Offset         int    `query:"offset"`
	IdentityNumber string `query:"identityNumber"`
	Name           string `query:"name"`
	Gender         string `query:"gender"`
	DepartmentID   string `query:"departmentId"` // Menggunakan string karena departmentId adalah string
}

// Normalize untuk set default values dan validasi
func (f *EmployeeFilter) Normalize() {
	if f.Limit <= 0 {
		f.Limit = 5
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}
