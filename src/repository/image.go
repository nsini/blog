package repository

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

type Image struct {
	gorm.Model
	ID        int64       `gorm:"column:id;primary_key" json:"id"`
	ImageName string      `gorm:"column:image_name" json:"image_name"`
	Extension null.String `gorm:"column:extension" json:"extension"`
	ImagePath null.String `gorm:"column:image_path" json:"image_path"`
	RealPath  null.String `gorm:"column:real_path" json:"real_path"`
	//ImageTime          null.Time   `gorm:"column:image_time" json:"image_time"`
	ImageStatus        null.Int    `gorm:"column:image_status" json:"image_status"`
	ImageSize          null.String `gorm:"column:image_size" json:"image_size"`
	Md5                null.String `gorm:"column:md5" json:"md5"`
	ClientOriginalMame null.String `gorm:"column:client_original_mame" json:"client_original_mame"`
	PostID             int64       `gorm:"column:post_id" json:"post_id"`
}

//var (
//ImageNotFound = errors.New("Image not found!")
//)

func (p *Image) TableName() string {
	return "images"
}

type ImageRepository interface {
	FindByPostIdLast(postId int64) (res *Image, err error)
	FindByPostIds(ids []int64) (res []*Image, err error)
	AddImage(img *Image) error
	ExistsImageByMd5(val string) bool
}

type image struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &image{db: db}
}

func (c *image) FindByPostIdLast(postId int64) (res *Image, err error) {
	var i Image
	if err = c.db.Last(&i, "post_id=?", postId).Error; err != nil {
		return &i, nil
	}
	return
}

func (c *image) FindByPostIds(ids []int64) (res []*Image, err error) {
	if err = c.db.Raw("SELECT image_name,image_path,MAX(id) id,post_id,real_path FROM `images`  WHERE `images`.`deleted_at` IS NULL AND ((post_id in (?))) GROUP BY post_id ORDER BY image_time DESC", ids).
		Scan(&res).Error; err != nil {
		return
	}
	return
}

func (c *image) AddImage(img *Image) error {
	return c.db.Save(img).Error
}

func (c *image) ExistsImageByMd5(val string) bool {
	var img Image
	if err := c.db.Where("md5 = ?", val).First(&img).Error; err != nil {
		return false
	}
	if img.Md5.String != "" {
		return true
	}
	return false
}
