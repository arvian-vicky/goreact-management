package repositories

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
	"fmt"
	"strings"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindAllPaginate(filter string, sort string, limit, offset int) ([]models.User, int64, error)
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := config.DB.Find(&users).Error
	return users, err
}

func (r *userRepository) FindAllPaginate(filter string, sort string, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := config.DB.Model(&models.User{})

	// Filtering: cari di nama atau email (sesuai kebutuhan)
	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ?", filterPattern, filterPattern)
	}

	// Count total sebelum limit & offset
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	fmt.Println("ini adalah sort: " + sort)
	// Sorting
	if sort != "" {
		// Misalnya sort=name atau sort=-name untuk descending
		if sort == "-id" {
			sort = "-internal_id"
		} else if sort == "id" {
			sort = "internal_id"
		}

		if strings.HasPrefix(sort, "-") {
			sort = strings.TrimPrefix(sort, "-") + " DESC"
		} else {
			sort += " ASC"
		}
		fmt.Println(sort)

		fmt.Println(sort)
		db = db.Order(sort)
	}

	// Pagination
	err := db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByPublicID(publicID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("public_id = ?", publicID).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return config.DB.Model(&models.User{}).
		Where("public_id = ?", user.PublicID).Updates(map[string]interface{}{
		"name": user.Name,
		//"email": newEmail,
	}).Error
}

func (r *userRepository) Delete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}
