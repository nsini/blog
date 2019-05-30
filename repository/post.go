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
	//User        User        `gorm:"foreignkey:UserRefer"`
}

var (
	PostNotFound = errors.New("post not found!")
)

func (p *Post) TableName() string {
	return "posts"
}

type PostRepository interface {
	Find(id int64) (res *Post, err error)
	FindBy(order, by string, limit, pageSize, offset int) ([]*Post, uint64, error)
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

	var user User

	if err = c.db.Find(&p, id).Related(&user, "ID").Error; err != nil {
		//if err = c.db.Find(&p, id).Error; err != nil {
		return nil, PostNotFound
	}

	//p.User = user
	return &p, nil
}

func (c *post) FindBy(order, by string, limit, pageSize, offset int) ([]*Post, uint64, error) {

	posts := make([]*Post, 0)
	var count uint64
	if err := c.db.Order(gorm.Expr(by + " " + order)).Where("push_time IS NOT NULL").Offset(offset).Limit(limit).Find(&posts).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}
