package repository

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"time"

	"gopkg.in/guregu/null.v3"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

type Comment struct {
	gorm.Model
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	PostID    null.Int    `gorm:"column:post_id" json:"post_id"`
	LogID     int64       `gorm:"column:log_id" json:"log_id"`
	APIUserID int         `gorm:"column:api_user_id" json:"api_user_id"`
	APIAction string      `gorm:"column:api_action" json:"api_action"`
	APIPostID null.Int    `gorm:"column:api_post_id" json:"api_post_id"`
	ThreadID  null.Int    `gorm:"column:thread_id" json:"thread_id"`
	ThreadKey null.String `gorm:"column:thread_key" json:"thread_key"`
	CommentIP null.String `gorm:"column:comment_ip" json:"comment_ip"`
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
	Message   null.String `gorm:"column:message" json:"message"`
	Status    null.String `gorm:"column:status" json:"status"`
	ParentID  null.Int    `gorm:"column:parent_id" json:"parent_id"`
	Type      null.Int    `gorm:"column:type" json:"type"`
	Agent     null.String `gorm:"column:agent" json:"agent"`
}

type comment struct {
	db *gorm.DB
}

func (p *comment) TableName() string {
	return "comment"
}

type CommentRepository interface {
	FindByPostId(postId int64) (res *Comment, err error)
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &comment{db: db}
}

func (c *comment) FindByPostId(postId int64) (res *Comment, err error) {
	var i Comment
	if err = c.db.Last(&i, "post_id=?", postId).Error; err != nil {
		return
	}
	return
}
