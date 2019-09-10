/**
 * @Time : 2019-09-10 10:57
 * @Author : solacowa@gmail.com
 * @File : meta
 * @Software: GoLand
 */

package types

type Tag struct {
	Id          int64  `gorm:"column:id;primary_key" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Slug        string `gorm:"column:slug" json:"slug"`
	Description string `gorm:"column:description" json:"description"`
	Count       int64  `gorm:"column:count;default(1)" json:"count"`
	Order       int    `gorm:"column:order;default(0)" json:"order"`
	Parent      int64  `gorm:"column:parent;default(0)" json:"parent"`
}

func (p *Tag) TableName() string {
	return "tags"
}
