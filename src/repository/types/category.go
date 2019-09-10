/**
 * @Time : 2019-09-10 11:10
 * @Author : solacowa@gmail.com
 * @File : category
 * @Software: GoLand
 */

package types

type Category struct {
	Id          int64  `gorm:"column:id;primary_key" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	ParentId    int64  `gorm:"column:parent_id;default(0)" json:"parent_id"`
}

func (p *Category) TableName() string {
	return "categories"
}
