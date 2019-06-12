package api

import (
	"context"
	"encoding/xml"
	"github.com/go-kit/kit/endpoint"
)

type postRequest struct {
	XMLName    xml.Name `xml:"methodCall"`
	MethodName string   `xml:"methodName"`
	Params     []params `xml:"params"`
}

type params struct {
	Param []value `xml:"param"`
}

type value struct {
	Value []vString `xml:"value"`
}

type vString struct {
	String string `xml:"string"`
}

type postResponse struct {
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)

		// todo 进行token验证

		rs, err := s.Post(ctx, req)
		return postResponse{rs, err}, err
	}
}
