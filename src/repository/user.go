package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/nsini/blog/src/repository/types"
)

type UserRepository interface {
	Find(username string) (res *types.User, err error)
}

type user struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &user{db: db}
}

func (c *user) Find(username string) (res *types.User, err error) {
	var rs types.User
	err = c.db.First(&rs, "username = ?", username).Error
	return &rs, err
}
