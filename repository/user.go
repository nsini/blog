package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	gorm.Model
	CreatedAt time.Time `gorm:"column:created_at"`
	Email     string    `gorm:"column:email"`
	ID        int       `gorm:"column:id;primary_key"`
	Password  string    `gorm:"column:password"`
	Username  string    `gorm:"column:username"`

	//Posts []Post `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"`
}

var (
	UserNotFound = errors.New("user not found!")
)

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
