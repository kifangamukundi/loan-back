package seeds

import (
	"log"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) error {
	seedFunctions := []struct {
		Name string
		Func func(*gorm.DB) error
	}{
		{Name: "Permissions", Func: SeedPermissions},
		{Name: "Roles", Func: SeedRolesAndAssignPermissions},
		{Name: "Users", Func: SeedUsers},
	}

	for _, seed := range seedFunctions {
		log.Printf("Seeding: %s", seed.Name)
		if err := seed.Func(db); err != nil {
			return err
		}
	}

	log.Println("All seeds ran successfully.")
	return nil
}
