package post

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/repository"
	"github.com/pkg/errors"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Detail(ctx context.Context, id int64) (rs *repository.Post, err error)
	List(ctx context.Context, order, by string, limit, pageSize, offset int) (rs []map[string]interface{}, count uint64, err error)
}

type service struct {
	post   repository.PostRepository
	user   repository.User
	logger log.Logger
}

/**
 * @Title 详情页
 */
func (c *service) Detail(ctx context.Context, id int64) (rs *repository.Post, err error) {
	detail, err := c.post.Find(id)
	if err != nil {
		return
	}

	if detail == nil {
		return nil, repository.PostNotFound
	}

	return detail, nil
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
			"id":         val.ID,
			"title":      val.Title,
			"desc":       val.Description,
			"publish_at": val.PushTime.Time.Format("2006/01/02 15:04:05"),
		})
	}

	return
}

func NewService(logger log.Logger, post repository.PostRepository, user repository.User) Service {
	return &service{
		post:   post,
		user:   user,
		logger: logger,
	}
}
