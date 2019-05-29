package post

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/blog/app/repository"
)

type postRequest struct {
	Id int64
}

type postResponse struct {
	Data *repository.Post `json:"data,omitempty"`
	Err  error            `json:"error,omitempty"`
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		rs, err := s.Detail(ctx, req.Id)
		return postResponse{rs, err}, err
	}
}
