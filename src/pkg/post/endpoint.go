package post

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type popularRequest struct {
}

type popularResponse struct {
	Data []map[string]interface{} `json:"data,omitempty"`
	Err  error                    `json:"error,omitempty"`
}

type postRequest struct {
	Id int64
}

type listRequest struct {
	order, by                string
	action, pageSize, offset int
}

type postResponse struct {
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error                  `json:"error,omitempty"`
}

type paginator struct {
	By       string `json:"by,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

type listResponse struct {
	Data      map[string]interface{} `json:"data,omitempty"`
	Count     int64                  `json:"count,omitempty"`
	Paginator paginator              `json:"paginator,omitempty"`
	Err       error                  `json:"error,omitempty"`
}

func (r listResponse) error() error { return r.Err }

type endpoints struct {
	GetEndpoint endpoint.Endpoint
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		rs, err := s.Get(ctx, req.Id)
		return postResponse{rs, err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		rs, count, other, err := s.List(ctx, req.order, req.by, req.action, req.pageSize, req.offset)
		return listResponse{
			Data: map[string]interface{}{
				"post":  rs,
				"other": other,
			},
			Count: count,
			Paginator: paginator{
				By:       req.by,
				Offset:   req.offset,
				PageSize: req.pageSize,
			},
			Err: err,
		}, err
	}
}

func makePopularEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := request.(popularRequest)
		rs, err := s.Popular(ctx)
		return popularResponse{rs, err}, err
	}
}
