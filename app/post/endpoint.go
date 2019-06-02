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
	order, by        string
	pageSize, offset int
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
	Data      []map[string]interface{} `json:"data,omitempty"`
	Count     uint64                   `json:"count,omitempty"`
	Paginator paginator                `json:"paginator,omitempty"`
	Err       error                    `json:"error,omitempty"`
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		rs, err := s.Detail(ctx, req.Id)
		return postResponse{rs, err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		rs, count, err := s.List(ctx, req.order, req.by, req.pageSize, req.offset)
		pager := paginator{
			By:       req.by,
			Offset:   req.offset,
			PageSize: req.pageSize,
		}
		return listResponse{rs, count, pager, err}, err
	}
}

func makePopularEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := request.(popularRequest)
		rs, err := s.Popular(ctx)
		return popularResponse{rs, err}, err
	}
}
