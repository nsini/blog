package api

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Post(ctx context.Context, method PostMethod, req postRequest) (rs newPostResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "post",
			"action", method,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, method, req)
}
