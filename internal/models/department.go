package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Department struct {
	ID           uint           `gorm:"primaryKey" json:"-"`                                                               // ID untuk auto increment
	DepartmentID string         `gorm:"size:10;not null;uniqueIndex:idx_department_id,where:deleted_at IS NULL" json:"id"` // Format: DEP-XX
	UserID       uint           `gorm:"not null" json:"-"`                                                                 // FK ke User
	Name         string         `gorm:"size:33;not null" json:"name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Employees []Employee `gorm:"foreignKey:DepartmentID;references:DepartmentID" json:"-"`
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

// Request structs
type CreateDepartmentRequest struct {
	Name string `json:"name" validate:"required,min=4,max=33"`
}

type UpdateDepartmentRequest struct {
	Name string `json:"name" validate:"required,min=4,max=33"`
}

// Hook untuk generate DepartmentID
func (d *Department) BeforeCreate(tx *gorm.DB) error {
	// Generate ID baru
	var lastID uint
	err := tx.Model(&Department{}).
		Order("id DESC").
		Limit(1).
		Pluck("id", &lastID).Error
	if err != nil {
		return err
	}

	// ID baru adalah lastID + 1
	newID := lastID + 1
	d.DepartmentID = fmt.Sprintf("DEP-%02d", newID)

	return nil
}

func (d *Department) ToResponse() *DepartmentResponse {
	return &DepartmentResponse{
		DepartmentID: d.DepartmentID,
		Name:         d.Name,
	}
}

func (f *DepartmentFilter) Normalize() {
	if f.Limit <= 0 {
		f.Limit = 5
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}
