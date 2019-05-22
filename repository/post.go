package repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Post struct {
	gorm.Model
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	IsMarkdown  bool      `json:"is_markdown"`
	ReadNum     int64     `json:"read_num"`
	Reviews     int64     `json:"reviews"`
	PushTime    time.Time `json:"push_time"`
	CreatedAt   time.Time `json:"created_at"`
	Action      int       `json:"action"`
	Star        int       `json:"star"`
	User        *User     `json:"user"`
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
	err = c.db.First(&res, id).Error
	if err != nil {
		return
	}
	return
}
