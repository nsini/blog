package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

type User struct {
	gorm.Model
	ID                  int         `gorm:"column:id;primary_key" json:"id"`
	Username            string      `gorm:"column:username" json:"username"`
	UsernameCanonical   string      `gorm:"column:username_canonical" json:"username_canonical"`
	Email               string      `gorm:"column:email" json:"email"`
	EmailCanonical      string      `gorm:"column:email_canonical" json:"email_canonical"`
	Enabled             int         `gorm:"column:enabled" json:"enabled"`
	Salt                string      `gorm:"column:salt" json:"salt"`
	Password            string      `gorm:"column:password" json:"password"`
	LastLogin           null.Time   `gorm:"column:last_login" json:"last_login"`
	Locked              int         `gorm:"column:locked" json:"locked"`
	Expired             int         `gorm:"column:expired" json:"expired"`
	ExpiresAt           null.Time   `gorm:"column:expires_at" json:"expires_at"`
	ConfirmationToken   null.String `gorm:"column:confirmation_token" json:"confirmation_token"`
	PasswordRequestedAt null.Time   `gorm:"column:password_requested_at" json:"password_requested_at"`
	Roles               string      `gorm:"column:roles" json:"roles"`
	CredentialsExpired  int         `gorm:"column:credentials_expired" json:"credentials_expired"`
	CredentialsExpireAt null.Time   `gorm:"column:credentials_expire_at" json:"credentials_expire_at"`
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
