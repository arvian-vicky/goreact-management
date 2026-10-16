package services

import (
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/repositories"
	"arfian-nurdavicky/project-management/utils"
	"errors"

	"github.com/google/uuid"
)

type UserService interface {
	GetAll() ([]models.User, error)
	GetAllPaginate(filter, sort string, limit, offset int) ([]models.User, int64, error)
	GetByID(id uint) (*models.User, error)
	GetByPublicID(id string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetAllPaginate(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	return s.repo.FindAllPaginate(filter, sort, limit, offset)
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}
func (s *userService) GetByPublicID(id string) (*models.User, error) {
	return s.repo.FindByPublicID(id)
}

func (s *userService) Create(user *models.User) error {
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.repo.Create(user)
}

func (s *userService) Update(user *models.User) error {
	return s.repo.Update(user)
}

func (s *userService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *userService) Register(user *models.User) error {
	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("email already registered")
	}
	hashed, err := utils.HashPassword(user.Password)

	if err != nil {
		return err
	}
	user.Password = hashed
	user.Role = "user"
	user.PublicID = uuid.New() // Atau uuid.UUID{}
	return s.repo.Create(user)
}

func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
