package repository

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"column:email"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	//ID       int64  `gorm:"column:id;primary_key"`
	//Posts []Post `gorm:"foreignkey:UserId"`
	Posts []Post `gorm:"polymorphic:Owner"`
}

func (p *User) TableName() string {
	return "users"
}

type UserRepository interface {
	Find(username string) (res *User, err error)
}

type user struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &user{db: db}
}

func (c *user) Find(username string) (res *User, err error) {
	return
}
