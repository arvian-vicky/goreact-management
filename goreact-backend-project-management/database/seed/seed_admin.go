package seed

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/utils"
	"log"
)

func SeedAdmin() {
	password, _ := utils.HashPassword("admin123")

	admin := models.User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: password,
		Role:     "admin",
	}

	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error; err != nil {
		log.Println("Failed to seed admin:", err)
	} else {
		log.Println("✅ Admin user seeded")
	}
}
