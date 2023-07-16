package service

import (
	"errors"

	"github.com/regmarcmem/gin-session-demo/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Signup(email string, password string) (*model.User, error) {
	user := model.User{}
	s.db.Where("email = ?", email).First(&user)
	if user.ID != 0 {
		err := errors.New("user already exists")
		return nil, err
	}

	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("failed to encrypt password")
		return nil, err
	}
	user = model.User{Email: email, Password: string(p)}

	s.db.Create(&user)
	return &user, nil
}
