/**
 * @Time : 2019-09-06 11:45
 * @Author : solacowa@gmail.com
 * @File : post
 * @Software: GoLand
 */

package types

import (
	"time"
)

type Post struct {
	Action      int       `gorm:"column:action" json:"action"`
	Content     string    `gorm:"column:content;type:text" json:"content"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
	Description string    `gorm:"column:description" json:"description"`
	Slug        string    `gorm:"column:slug" json:"slug"`
	ID          int64     `gorm:"column:id;primary_key" json:"id"`
	IsMarkdown  bool      `gorm:"column:is_markdown" json:"is_markdown"`
	PushTime    time.Time `gorm:"column:push_time" json:"push_time"`
	ReadNum     int64     `gorm:"column:read_num" json:"read_num"`
	Reviews     int64     `gorm:"column:reviews" json:"reviews"`
	Star        int       `gorm:"column:star" json:"star"`
	Status      int       `gorm:"column:status" json:"status"`
	Title       string    `gorm:"column:title" json:"title"`
	UserID      int64     `gorm:"column:user_id" json:"user_id"`
	User        User      `gorm:"ForeignKey:id;AssociationForeignKey:user_id"`
}

func (p *Post) TableName() string {
	return "posts"
}
