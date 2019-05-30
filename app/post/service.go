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
	List(ctx context.Context) (rs map[string]interface{}, err error)
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
func (c *service) List(ctx context.Context) (rs map[string]interface{}, err error) {

	return
}

func NewService(logger log.Logger, post repository.PostRepository, user repository.User) Service {
	return &service{
		post:   post,
		user:   user,
		logger: logger,
	}
}
