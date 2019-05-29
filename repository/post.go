package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/guregu/null.v3"
	"time"
)

type Post struct {
	gorm.Model
	Action      int         `gorm:"column:action"`
	Content     string      `gorm:"column:content"`
	CreatedAt   time.Time   `gorm:"column:created_at"`
	Description null.String `gorm:"column:description"`
	ID          int64       `gorm:"column:id;primary_key"`
	IsMarkdown  null.Int    `gorm:"column:is_markdown"`
	PushTime    null.Time   `gorm:"column:push_time"`
	ReadNum     int64       `gorm:"column:read_num"`
	Reviews     int64       `gorm:"column:reviews"`
	Star        null.Int    `gorm:"column:star"`
	Status      int         `gorm:"column:status"`
	Title       string      `gorm:"column:title"`
	UserID      null.Int    `gorm:"column:user_id"`
}

var (
	PostNotFound = errors.New("post not found!")
)

func (p *Post) TableName() string {
	return "posts"
}

type PostRepository interface {
	Find(id int64) (res *Post, err error)
}

type post struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &post{db: db}
}

func (c *post) Find(id int64) (res *Post, err error) {
	p := Post{
		ID: id,
	}
	if err = c.db.Find(&p, id).Error; err != nil {
		return nil, PostNotFound
	}
	return &p, nil
}
