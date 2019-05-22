package repository

import (
	"time"
)

type Post struct {
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
}

func NewPostRepository() PostRepository {

	return &post{}
}

func (c *post) Find(id int64) (res *Post, err error) {

	return &Post{
		Id:          1,
		Title:       "hello",
		Description: "world",
		Content:     "none",
		CreatedAt:   time.Now(),
	}, nil
}
