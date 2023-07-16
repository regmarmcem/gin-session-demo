package service

import (
	"github.com/regmarcmem/gin-session-demo/model"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) Signup(email string, password string) (*model.User, error) {
	return &model.User{Email: email, Password: password}, nil
}
