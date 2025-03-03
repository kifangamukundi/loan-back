package migrations

import (
	"log"

	"github.com/kifangamukundi/gm/loan/database"
	"github.com/kifangamukundi/gm/loan/models"
)

// RunMigrations runs migrations for all models
func RunMigrations() {
	err := database.DB.AutoMigrate(
		&models.Country{},
		&models.Region{},
		&models.County{},
		&models.SubCounty{},
		&models.Ward{},
		&models.Location{},
		&models.SubLocation{},
		&models.Village{},
		&models.Road{},
		&models.Plot{},
		&models.Unit{},
		&models.Role{},
		&models.Permission{},

		// Dependent tables
		&models.User{},
		&models.Agent{},
		&models.Group{},
		&models.Member{},
		&models.Loan{},
		&models.Officer{},
		&models.Disbursement{},
		&models.Payment{},

		// Join tables and associations
		&models.RolePermission{},
		&models.UserRole{},
		&models.GroupMember{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		log.Println("Database migration completed successfully!")
	}
}
