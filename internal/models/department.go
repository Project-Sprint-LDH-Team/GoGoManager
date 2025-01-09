package models

import (
	"time"
)

type Department struct {
	ID        string     `gorm:"primaryKey;type:varchar(36)" json:"department_id"`
	UserID    uint       `gorm:"not null;index" json:"-"`
	Name      string     `gorm:"size:33;not null" json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`

	// Relations
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Employees []Employee `gorm:"foreignKey:DepartmentID" json:"-"`
}

type CreateDepartmentRequest struct {
	Name string `json:"name" validate:"required,min=4,max=33"`
}

type UpdateDepartmentRequest struct {
	Name string `json:"name" validate:"required,min=4,max=33"`
}

type DepartmentResponse struct {
	DepartmentID string `json:"departmentId"`
	Name         string `json:"name"`
}

type DepartmentFilter struct {
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
	Name   string `query:"name"`
}

func (f *DepartmentFilter) Normalize() {
	if f.Limit <= 0 {
		f.Limit = 5
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}

func (d *Department) ToResponse() *DepartmentResponse {
	return &DepartmentResponse{
		DepartmentID: d.ID,
		Name:         d.Name,
	}
}
