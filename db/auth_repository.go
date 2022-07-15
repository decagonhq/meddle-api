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
	FindUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	AddToBlackList(blacklist *models.BlackList) error
	TokenInBlacklist(token *string) bool
}

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *GormDB) AuthRepository {
	return &AuthRepo{db.DB}
}

func (a *AuthRepo) CreateUser(user *models.User) (*models.User, error) {
	db := a.DB
	err := db.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}
	return user, nil
}

func (a *AuthRepo) FindUserByUsername(username string) (*models.User, error) {
	db := a.DB
	user := &models.User{}
	err := db.Where("email = ? OR username = ?", username, username).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}

func (a *AuthRepo) FindUserByEmail(email string) (*models.User, error) {
	var user *models.User
	userEmail := a.DB.Where("email = ?", email).First(&user)
	return user, userEmail.Error
}

func (a *AuthRepo) UpdateUser(user *models.User) error {
	return nil
}

func (a *AuthRepo) AddToBlackList(blacklist *models.BlackList) error {
	return nil
}

func (a *AuthRepo) TokenInBlacklist(token *string) bool {
	return false
}
