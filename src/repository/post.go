package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
	"time"
)

type Post struct {
	//gorm.Model
	Action      int         `gorm:"column:action"`
	Content     string      `gorm:"column:content"`
	CreatedAt   time.Time   `gorm:"column:created_at"`
	UpdatedAt   time.Time   `gorm:"column:updated_at"`
	Description null.String `gorm:"column:description"`
	Slug        null.String `gorm:"column:slug"`
	ID          int64       `gorm:"column:id;primary_key"`
	IsMarkdown  null.Int    `gorm:"column:is_markdown"`
	PushTime    null.Time   `gorm:"column:push_time"`
	ReadNum     int64       `gorm:"column:read_num"`
	Reviews     int64       `gorm:"column:reviews"`
	Star        null.Int    `gorm:"column:star"`
	Status      int         `gorm:"column:status"`
	Title       string      `gorm:"column:title"`
	UserID      null.Int    `gorm:"column:user_id"`
	User        User        `gorm:"ForeignKey:id;AssociationForeignKey:user_id"`
}

var (
	PostNotFound = errors.New("post not found!")
)

func (p *Post) TableName() string {
	return "posts"
}

type PostRepository interface {
	Find(id int64) (res *Post, err error)
	FindBy(action int, order, by string, pageSize, offset int) ([]*Post, int64, error)
	Popular() (posts []*Post, err error)
	SetReadNum(p *Post) error
	Create(p *Post) error
}

type post struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &post{db: db}
}

func (c *post) Find(id int64) (res *Post, err error) {
	var p Post

	if err = c.db.Preload("User").
		Find(&p, "id = ?", id).Error; err != nil {
		return nil, PostNotFound
	}
	return &p, nil
}

func (c *post) FindBy(action int, order, by string, pageSize, offset int) ([]*Post, int64, error) {
	posts := make([]*Post, 0)
	var count int64
	if err := c.db.Model(&posts).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,username")
	}).
		Where("action = ?", action).
		Order(gorm.Expr(by + " " + order)).
		Where("push_time IS NOT NULL").Count(&count).
		Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

func (c *post) Popular() (posts []*Post, err error) {
	if err = c.db.Order("read_num DESC").Limit(5).Find(&posts).Error; err != nil {
		return
	}
	return
}

func (c *post) SetReadNum(p *Post) error {
	p.ReadNum += 1
	return c.db.Exec("UPDATE `posts` SET `read_num` = ?  WHERE `posts`.`deleted_at` IS NULL AND `posts`.`id` = ?", p.ReadNum, p.ID).Error
}

func (c *post) Create(p *Post) error {
	if err := c.db.Save(p).Error; err != nil {
		return err
	}
	return nil
}
