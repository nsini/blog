/**
 * @Time : 2019-09-10 11:05
 * @Author : solacowa@gmail.com
 * @File : meta
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/nsini/blog/src/repository/types"
)

type MetaType string

const (
	MetaCategory MetaType = "category"
	MetaTag      MetaType = "tag"
)

type TagRepository interface {
	FirstOrCreate(name string) (meta *types.Tag, err error)
	List(limit int) (metas []*types.Tag, err error)
}

type tag struct {
	db *gorm.DB
}

func (c *tag) List(limit int) (metas []*types.Tag, err error) {
	err = c.db.Model(&types.Tag{}).Order("id desc").Limit(limit).Find(&metas).Error
	return
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tag{db: db}
}

func (c *tag) FirstOrCreate(name string) (tag *types.Tag, err error) {
	t := types.Tag{
		Name:        name,
		Description: name,
	}
	err = c.db.FirstOrCreate(&t, types.Tag{
		Name: name,
	}).Error

	return &t, err
}
