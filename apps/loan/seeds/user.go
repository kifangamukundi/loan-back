package seeds

import (
	"log"

	"github.com/kifangamukundi/gm/libs/auths"
	"github.com/kifangamukundi/gm/loan/database"
	"github.com/kifangamukundi/gm/loan/models"
	"github.com/kifangamukundi/gm/libs/repositories"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	repo := repositories.NewGormRepository(database.DB)

	users := []struct {
		User  models.User
		Roles []string
	}{
		{
			User: models.User{
				FirstName:    "Kifanga",
				LastName:     "Mukundi",
				Email:        "pkmymcmbnhs@gmail.com",
				MobileNumber: "+254702817040",
				IsActive:     true,
			},
			Roles: []string{"Admin"},
		},
	}

	for _, userInfo := range users {
		var user models.User

		if err := db.Where("email = ?", userInfo.User.Email).First(&user).Error; err == nil {
			log.Printf("User %s already exists, skipping seed.", userInfo.User.Email)
			continue
		} else if err != gorm.ErrRecordNotFound {
			log.Printf("Error fetching user by email %s: %v", userInfo.User.Email, err)
			return err
		}

		hashedPassword, err := auths.HashPassword("123456789")
		if err != nil {
			log.Println("Password Hashing Failed")
			return err
		}

		userInfo.User.Password = hashedPassword
		if err := repo.Create(&userInfo.User); err != nil {
			log.Printf("Error creating user %s: %v", userInfo.User.Email, err)
			return err
		}

		var roles []models.Role
		if err := database.DB.Where("role_name IN ?", userInfo.Roles).Find(&roles).Error; err != nil {
			log.Printf("Error fetching roles for user %s: %v", userInfo.User.Email, err)
			return err
		}

		if err := database.DB.Model(&userInfo.User).Association("Roles").Replace(roles); err != nil {
			log.Printf("Error associating roles for user %s: %v", userInfo.User.Email, err)
			return err
		}

		log.Printf("User %s seeded successfully with roles.", userInfo.User.Email)
	}

	return nil
}
