package post

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/nsini/blog/repository"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) Detail(ctx context.Context, id int64) (rs *repository.Post, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "detail").Add(1)
		s.requestLatency.With("method", "detail").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Detail(ctx, id)
}

func (s *instrumentingService) List(ctx context.Context) (rs map[string]interface{}, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list").Add(1)
		s.requestLatency.With("method", "list").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.List(ctx)
}