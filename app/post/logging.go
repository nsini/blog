package post

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/repository"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Detail(ctx context.Context, id int64) (rs *repository.Post, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "detail",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, id)
}

func (s *loggingService) List(ctx context.Context) (rs map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx)
}
