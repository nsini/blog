package post

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/repository"
	"github.com/pkg/errors"
	"strconv"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Detail(ctx context.Context, id int64) (rs map[string]interface{}, err error)
	List(ctx context.Context, order, by string, limit, pageSize, offset int) (rs []map[string]interface{}, count uint64, err error)
	Popular(ctx context.Context) (rs []map[string]interface{}, err error)
}

type service struct {
	post   repository.PostRepository
	user   repository.UserRepository
	image  repository.ImageRepository
	logger log.Logger
}

/**
 * @Title 详情页
 */
func (c *service) Detail(ctx context.Context, id int64) (rs map[string]interface{}, err error) {
	detail, err := c.post.Find(id)
	if err != nil {
		return
	}

	if detail == nil {
		return nil, repository.PostNotFound
	}

	var headerImage string

	if image, err := c.image.FindByPostIdLast(id); err == nil && image != nil {
		headerImage = "/image/" + image.ImagePath.String
	}

	return map[string]interface{}{
		"content":      detail.Content,
		"title":        detail.Title,
		"publish_at":   detail.PushTime.Time.Format("2006/01/02 15:04:05"),
		"updated_at":   detail.UpdatedAt,
		"author":       detail.User.Username,
		"comment":      4,
		"banner_image": headerImage,
	}, nil
}

/**
 * @Title 列表页
 */
func (c *service) List(ctx context.Context, order, by string, limit, pageSize, offset int) (rs []map[string]interface{}, count uint64, err error) {
	// 取列表 判断搜索、分类、Tag条件
	// 取最多阅读

	posts, count, err := c.post.FindBy(order, by, limit, pageSize, offset)
	if err != nil {
		return
	}

	_ = c.logger.Log("count", count)

	for _, val := range posts {
		rs = append(rs, map[string]interface{}{
			"id":         strconv.FormatUint(uint64(val.Model.ID), 10),
			"title":      val.Title,
			"desc":       val.Description,
			"publish_at": val.PushTime.Time.Format("2006/01/02 15:04:05"),
		})
	}

	return
}

func (c *service) Popular(ctx context.Context) (rs []map[string]interface{}, err error) {

	posts, err := c.post.Popular()
	if err != nil {
		return
	}

	for _, post := range posts {
		fmt.Println(post.Title)
	}

	return
}

func NewService(logger log.Logger, post repository.PostRepository, user repository.UserRepository, image repository.ImageRepository) Service {
	return &service{
		post:   post,
		user:   user,
		image:  image,
		logger: logger,
	}
}
