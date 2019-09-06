/**
 * @Time : 2019-09-06 11:32
 * @Author : solacowa@gmail.com
 * @File : init
 * @Software: GoLand
 */

package repository

import "github.com/jinzhu/gorm"

type Repository interface {
	User() UserRepository
	Post() PostRepository
	Image() ImageRepository
	Comment() CommentRepository
}

type store struct {
	db      *gorm.DB
	user    UserRepository
	post    PostRepository
	image   ImageRepository
	comment CommentRepository
}

func NewRepository(db *gorm.DB) Repository {
	return &store{
		user:    NewUserRepository(db),
		post:    NewPostRepository(db),
		image:   NewImageRepository(db),
		comment: NewCommentRepository(db),
	}
}

func (c *store) User() UserRepository {
	return c.user
}

func (c *store) Post() PostRepository {
	return c.post
}

func (c *store) Image() ImageRepository {
	return c.image
}

func (c *store) Comment() CommentRepository {
	return c.comment
}
