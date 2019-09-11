package post

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/src/config"
	"github.com/nsini/blog/src/repository"
	"strconv"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Get(ctx context.Context, id int64) (rs map[string]interface{}, err error)
	List(ctx context.Context, order, by string, action, pageSize, offset int) (rs []map[string]interface{}, count int64, err error)
	Popular(ctx context.Context) (rs []map[string]interface{}, err error)
}

type service struct {
	repository repository.Repository
	logger     log.Logger
	config     *config.Config
}

/**
 * @Title 详情页
 */
func (c *service) Get(ctx context.Context, id int64) (rs map[string]interface{}, err error) {
	detail, err := c.repository.Post().Find(id)
	if err != nil {
		return
	}

	if detail == nil {
		return nil, repository.PostNotFound
	}

	defer func() {
		go func() {
			if err = c.repository.Post().SetReadNum(detail); err != nil {
				_ = c.logger.Log("post.SetReadNum", err.Error())
			}
		}()
	}()

	var headerImage string

	if image, err := c.repository.Image().FindByPostIdLast(id); err == nil && image != nil {
		headerImage = c.config.GetString("server", "image_domain") + "/" + image.ImagePath
	}

	return map[string]interface{}{
		"content":      detail.Content,
		"title":        detail.Title,
		"publish_at":   detail.PushTime.Format("2006/01/02 15:04:05"),
		"updated_at":   detail.UpdatedAt,
		"author":       detail.User.Username,
		"comment":      detail.Reviews,
		"banner_image": headerImage,
	}, nil
}

/**
 * @Title 列表页
 */
func (c *service) List(ctx context.Context, order, by string, action, pageSize, offset int) (rs []map[string]interface{}, count int64, err error) {
	// 取列表 判断搜索、分类、Tag条件
	// 取最多阅读

	posts, count, err := c.repository.Post().FindBy(action, order, by, pageSize, offset)
	if err != nil {
		return
	}

	var postIds []int64

	for _, post := range posts {
		postIds = append(postIds, post.ID)
	}

	images, err := c.repository.Image().FindByPostIds(postIds)
	if err == nil && images == nil {
		_ = c.logger.Log("c.image.FindByPostIds", "postIds", "err", err)
	}

	imageMap := make(map[int64]string, len(images))
	for _, image := range images {
		imageMap[image.PostID] = imageUrl(image.ImagePath, c.config.GetString("server", "image_domain"))
	}

	_ = c.logger.Log("count", count)

	for _, val := range posts {
		imageUrl, ok := imageMap[val.ID]
		if !ok {
			_ = c.logger.Log("postId", val.ID, "image", ok)
		}
		rs = append(rs, map[string]interface{}{
			"id":         strconv.FormatUint(uint64(val.ID), 10),
			"title":      val.Title,
			"desc":       val.Description,
			"publish_at": val.PushTime.Format("2006/01/02 15:04:05"),
			"image_url":  imageUrl,
			"comment":    val.Reviews,
			"author":     val.User.Username,
			"tags":       val.Tags,
		})
	}

	return
}

/**
 * @Title 受欢迎的
 */
func (c *service) Popular(ctx context.Context) (rs []map[string]interface{}, err error) {

	posts, err := c.repository.Post().Popular()
	if err != nil {
		return
	}

	var postIds []int64

	for _, post := range posts {
		postIds = append(postIds, post.ID)
	}

	images, err := c.repository.Image().FindByPostIds(postIds)
	if err == nil && images == nil {
		_ = c.logger.Log("c.image.FindByPostIds", "postIds", "err", err)
	}

	imageMap := make(map[int64]string, len(images))
	for _, image := range images {
		imageMap[image.PostID] = imageUrl(image.ImagePath, c.config.GetString("server", "image_domain"))
	}

	for _, post := range posts {
		imageUrl, ok := imageMap[post.ID]
		if !ok {
			_ = c.logger.Log("postId", post.ID, "image", ok)
		}

		desc := []rune(post.Description)
		rs = append(rs, map[string]interface{}{
			"title":     post.Title,
			"desc":      string(desc[:40]),
			"id":        post.ID,
			"image_url": imageUrl,
		})
	}

	return
}

func imageUrl(path, imageDomain string) string {
	return imageDomain + "/" + path
}

func NewService(logger log.Logger, cf *config.Config, repository repository.Repository) Service {
	return &service{
		repository: repository,
		logger:     logger,
		config:     cf,
	}
}
