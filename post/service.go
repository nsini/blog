package post

import (
	"context"
	"github.com/nsini/blog/repository"
	"github.com/pkg/errors"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Detail(ctx context.Context, id int64) (rs map[string]interface{}, err error)
}

type service struct {
	post repository.PostRepository
	user repository.User
}

func (c *service) Detail(ctx context.Context, id int64) (rs map[string]interface{}, err error) {

	detail, err := c.post.Find(id)
	if err != nil {
		return
	}

	return map[string]interface{}{
		"id":      detail.Id,
		"title":   detail.Title,
		"content": detail.Content,
	}, nil
}

func NewService(post repository.PostRepository, user repository.User) Service {
	return &service{
		post: post,
		user: user,
	}
}
