package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/nsini/blog/src/repository/types"
)

var (
	PostNotFound = errors.New("post not found!")
)

type PostRepository interface {
	Find(id int64) (res *types.Post, err error)
	FindBy(action int, order, by string, pageSize, offset int) ([]*types.Post, int64, error)
	Popular() (posts []*types.Post, err error)
	SetReadNum(p *types.Post) error
	Create(p *types.Post) error
	Update(p *types.Post) error
}

type post struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &post{db: db}
}

func (c *post) Update(p *types.Post) error {
	return c.db.Model(p).Where("id = ?", p.ID).Update(p).Error
}

func (c *post) Find(id int64) (res *types.Post, err error) {
	var p types.Post

	if err = c.db.Model(&p).
		Preload("User").
		Preload("Categories").
		Find(&p, "id = ?", id).Error; err != nil {
		return nil, PostNotFound
	}
	return &p, nil
}

func (c *post) FindBy(action int, order, by string, pageSize, offset int) ([]*types.Post, int64, error) {
	posts := make([]*types.Post, 0)
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

func (c *post) Popular() (posts []*types.Post, err error) {
	if err = c.db.Order("read_num DESC").Limit(5).Find(&posts).Error; err != nil {
		return
	}
	return
}

func (c *post) SetReadNum(p *types.Post) error {
	p.ReadNum += 1
	return c.db.Exec("UPDATE `posts` SET `read_num` = ?  WHERE `posts`.`deleted_at` IS NULL AND `posts`.`id` = ?", p.ReadNum, p.ID).Error
}

func (c *post) Create(p *types.Post) error {
	if err := c.db.Save(p).Error; err != nil {
		return err
	}
	return nil
}
