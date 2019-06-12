package api

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/config"
	"github.com/nsini/blog/repository"
)

type Service interface {
	Post(ctx context.Context, req postRequest) (rs map[string]interface{}, err error)
}

type service struct {
	post   repository.PostRepository
	user   repository.UserRepository
	image  repository.ImageRepository
	logger log.Logger
	config config.Config
}

func (c *service) Post(ctx context.Context, req postRequest) (rs map[string]interface{}, err error) {

	_ = c.logger.Log("methodName", req.MethodName)

	for _, v := range req.Params {
		for _, val := range v.Param {
			_ = c.logger.Log("string", val.Value[0].String)
		}
	}

	return
}

func NewService(logger log.Logger, cf config.Config, post repository.PostRepository, user repository.UserRepository, image repository.ImageRepository) Service {
	return &service{
		post:   post,
		user:   user,
		image:  image,
		logger: logger,
		config: cf,
	}
}
