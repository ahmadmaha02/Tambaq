package user

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]User, error)
	FindByID(ID int) (User, error)
	Register(user User) (User, error)
	RegisterMember(user User) (User, error)
	Login()([]User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Register(user User) (User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *repository) RegisterMember(user User) (User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) FindAll() ([]User, error) {
	var books []User
	err := r.db.Find(&books).Error
	return books, err
}