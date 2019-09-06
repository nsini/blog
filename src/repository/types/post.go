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
	Action      int       `gorm:"column:action"`
	Content     string    `gorm:"column:content"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Description string    `gorm:"column:description"`
	Slug        string    `gorm:"column:slug"`
	ID          int64     `gorm:"column:id;primary_key"`
	IsMarkdown  bool      `gorm:"column:is_markdown"`
	PushTime    time.Time `gorm:"column:push_time"`
	ReadNum     int64     `gorm:"column:read_num"`
	Reviews     int64     `gorm:"column:reviews"`
	Star        int       `gorm:"column:star"`
	Status      int       `gorm:"column:status"`
	Title       string    `gorm:"column:title"`
	UserID      int64     `gorm:"column:user_id"`
	User        User      `gorm:"ForeignKey:id;AssociationForeignKey:user_id"`
}

func (p *Post) TableName() string {
	return "posts"
}
