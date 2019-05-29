package about

import (
	"context"
	"github.com/go-kit/kit/metrics"
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

func (s *instrumentingService) Detail(ctx context.Context) (rs map[string]interface{}, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "about").Add(1)
		s.requestLatency.With("method", "about").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.About(ctx)
}
