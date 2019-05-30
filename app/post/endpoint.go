package post

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/blog/repository"
)

type postRequest struct {
	Id int64
}

type listRequest struct {
}

type postResponse struct {
	Data *repository.Post `json:"data,omitempty"`
	Err  error            `json:"error,omitempty"`
}

type listResponse struct {
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error                  `json:"error,omitempty"`
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
		//req := request.(listRequest)
		rs, err := s.List(ctx)
		return listResponse{rs, err}, err
	}
}
