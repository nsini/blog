package post

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type postRequest struct {
	Id int64
}

type bookCargoResponse struct {
	Res map[string]interface{} `json:"res,omitempty"`
	Err error                  `json:"error,omitempty"`
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		rs, err := s.Detail(ctx, req.Id)
		return bookCargoResponse{rs, err}, nil
	}
}
