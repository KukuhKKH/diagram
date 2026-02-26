package repository

import (
	"git.dev.siap.id/kukuhkkh/app-diagram/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-diagram/internal/bootstrap/database"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/helpers"
)

type userRepository struct {
	DB *database.Database
}

//go:generate mockgen -destination=article_repository_mock.go -package=repository . UserRepository
type UserRepository interface {
	FindUserByID(id uint64) (user *schema.User, err error)
	FindUserByEmail(email string) (user *schema.User, err error)
	FindUserByLogtoSub(sub string) (user *schema.User, err error)
	CheckUserByEmail(email string) (user *schema.User)
	CreateUser(user *schema.User) (res *schema.User, err error)
	UpdateUser(user *schema.User) error
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (_i *userRepository) FindUserByID(id uint64) (user *schema.User, err error) {
	if err := _i.DB.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return
}

func (_i *userRepository) FindUserByEmail(email string) (user *schema.User, err error) {
	if err := _i.DB.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return
}

func (_i *userRepository) FindUserByLogtoSub(sub string) (user *schema.User, err error) {
	if err := _i.DB.DB.Where("logto_sub = ?", sub).First(&user).Error; err != nil {
		return nil, err
	}
	return
}

func (_i *userRepository) CheckUserByEmail(email string) (user *schema.User) {
	_ = _i.DB.DB.Where("email = ?", email).First(&user).Error
	return
}

func (_i *userRepository) CreateUser(user *schema.User) (res *schema.User, err error) {
	if user.Password != nil {
		hashed := helpers.Hash([]byte(*user.Password))
		user.Password = &hashed
	}

	if err := _i.DB.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (_i *userRepository) UpdateUser(user *schema.User) error {
	return _i.DB.DB.Save(user).Error
}
