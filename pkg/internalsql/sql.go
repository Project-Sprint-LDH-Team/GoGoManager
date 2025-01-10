package internalsql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func Connect(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		// Optimalkan koneksi database
		PrepareStmt: true, // Cache prepared statements
		// Disable logger pada production
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("error connecting to database %+v\n", err)
		return nil, err
	}

	// Set connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)           // Jumlah minimum koneksi yang idle
	sqlDB.SetMaxOpenConns(100)          // Jumlah maksimum koneksi yang dibuka
	sqlDB.SetConnMaxLifetime(time.Hour) // Waktu maksimum sebuah koneksi dapat digunakan

	// Create Users table
	if err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            name VARCHAR(52),
            user_image_uri VARCHAR(255),
            company_name VARCHAR(52),
            company_image_uri VARCHAR(255),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `).Error; err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	// Create Departments table
	// Create Departments table
	if err := db.Exec(`
        CREATE TABLE IF NOT EXISTS departments (
            id SERIAL PRIMARY KEY,
            department_id VARCHAR(10) NOT NULL,
            user_id INTEGER NOT NULL,
            name VARCHAR(33) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP WITH TIME ZONE,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )
    `).Error; err != nil {
		return nil, fmt.Errorf("failed to create departments table: %w", err)
	}

	// Buat partial unique index dalam query terpisah
	if err := db.Exec(`
        CREATE UNIQUE INDEX IF NOT EXISTS idx_departments_department_id 
        ON departments (department_id) 
        WHERE deleted_at IS NULL
    `).Error; err != nil {
		return nil, fmt.Errorf("failed to create department unique index: %w", err)
	}

	// Create Employees table
	if err := db.Exec(`
        CREATE TABLE IF NOT EXISTS employees (
            id SERIAL PRIMARY KEY,
            department_id VARCHAR(10) NOT NULL,
            identity_number VARCHAR(33) UNIQUE NOT NULL,
            name VARCHAR(33) NOT NULL,
            employee_image_uri VARCHAR(255),
            gender VARCHAR(6) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP WITH TIME ZONE,
            FOREIGN KEY (department_id) REFERENCES departments(department_id)
        )
    `).Error; err != nil {
		return nil, fmt.Errorf("failed to create employees table: %w", err)
	}

	// Create Files table
	if err := db.Exec(`
        CREATE TABLE IF NOT EXISTS files (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL,
            filename VARCHAR(255) NOT NULL,
            file_uri VARCHAR(255) NOT NULL,
            file_type VARCHAR(50) NOT NULL,
            file_size BIGINT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )
    `).Error; err != nil {
		return nil, fmt.Errorf("failed to create files table: %w", err)
	}

	// Create indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_departments_user_id ON departments(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_departments_name ON departments(name)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_employees_department_id ON employees(department_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_employees_name ON employees(name)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_files_user_id ON files(user_id)")

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to verify database connection: %w", err)
	}

	return db, nil
}

// Utility function untuk close database connection dengan graceful
func CloseDatabaseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
