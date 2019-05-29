package post

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/app/repository"
	"github.com/pkg/errors"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Detail(ctx context.Context, id int64) (rs *repository.Post, err error)
}

type service struct {
	post   repository.PostRepository
	user   repository.User
	logger log.Logger
}

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

func NewService(logger log.Logger, post repository.PostRepository, user repository.User) Service {
	return &service{
		post:   post,
		user:   user,
		logger: logger,
	}
}
