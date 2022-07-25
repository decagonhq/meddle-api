package db

import (
	"fmt"

	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
)

// DB provides access to the different db
type AuthRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	AddToBlackList(blacklist *models.BlackList) error
	TokenInBlacklist(token string) bool
}

type authRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *GormDB) AuthRepository {
	return &authRepo{db.DB}
}

func (a *authRepo) CreateUser(user *models.User) (*models.User, error) {
	db := a.DB
	err := db.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}
	return user, nil
}

func (a *authRepo) FindUserByUsername(username string) (*models.User, error) {
	db := a.DB
	user := &models.User{}
	err := db.Where("email = ? OR username = ?", username, username).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}

func (a *authRepo) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := a.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepo) UpdateUser(user *models.User) error {
	return nil
}

func (a *authRepo) AddToBlackList(blacklist *models.BlackList) error {
	return nil
}

func (a *authRepo) TokenInBlacklist(token string) bool {
	result := a.DB.Where("token = ?", token).Find(&models.BlackList{})
	return result.Error != nil
}
