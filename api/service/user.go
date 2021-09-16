package service

import (
	"database/sql"
	"errors"
	"main/models"
	"main/repository"
	"main/utils"

	"github.com/google/uuid"
)

var (
	userRepo = repository.NewUserRepository()
)

type UserService struct {
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{}
}

func (us *UserService) Validate(u *models.User) error {
	if u.Email == "" || u.PasswordHash == "" {
		return errors.New("Invalid value for email or password\n")
	}
	return nil
}

func (us *UserService) CreateUser(u *models.User) (*models.User, error) {
	var err error
	u.ID = uuid.NewString()
	u.PasswordHash, err = utils.HashPassword(u.PasswordHash)
	if err != nil {
		return nil, err
	}
	u, err = userRepo.Create(u)
	return u, err
}
func (us *UserService) FindUser(u *models.User) (*models.User, error) {
	var err error
	u, err = userRepo.FindUser(u)
	return u, err
}
