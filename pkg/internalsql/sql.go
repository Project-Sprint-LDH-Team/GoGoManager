package internalsql

import (
	"fmt"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Connect(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting to database %+v\n", err)
		return nil, err
	}

	// Drop foreign key constraint jika ada
	db.Exec("ALTER TABLE IF EXISTS employees DROP CONSTRAINT IF EXISTS fk_departments_employees")

	// Drop kolom department_id di employees jika ada
	db.Exec("ALTER TABLE IF EXISTS employees DROP COLUMN IF EXISTS department_id")

	// Auto Migrate dalam urutan yang benar
	err = db.AutoMigrate(
		&models.User{},
		&models.Department{}, // Migrate Department dulu
		&models.Employee{},   // Baru Employee
		&models.File{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
