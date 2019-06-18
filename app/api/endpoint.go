package api

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"strconv"
)

type PostMethod string

const (
	PostCreate     PostMethod = "metaWeblog.newPost"
	GetCategories  PostMethod = "metaWeblog.getCategories"
	NewMediaObject PostMethod = "metaWeblog.newMediaObject"
	GetPost        PostMethod = "metaWeblog.getPost"
	GetUsersBlogs  PostMethod = "blogger.getUsersBlogs"
)

func (c PostMethod) String() string {
	return string(c)
}

var NoPermission = errors.New("not permission!")

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		var err error
		var resp interface{}

		// todo 进行token验证
		fmt.Println("username", req.Params.Param[1].Value.String)
		fmt.Println("password", req.Params.Param[2].Value.String)

		switch PostMethod(req.MethodName) {
		case GetUsersBlogs:
			// todo check response
			{
				resp = &getUsersBlogsResponse{
					Params: params{
						Param: param{
							Value: value{
								Array: array{
									Data: data{
										Value: dataValue{
											Struct: valStruct{
												Member: []member{
													{Name: "isAdmin", Value: memberValue{String: "1"}},
													{Name: "url", Value: memberValue{String: "http://localhost:8080"}},
													{Name: "blogid", Value: memberValue{String: "1"}},
													{Name: "blogName", Value: memberValue{String: "nsini"}},
												},
											},
										},
									},
								},
							},
						},
					},
				}
			}
		case GetCategories:
			return s.GetCategories(ctx, req)
		case PostCreate:
			rs, err := s.Post(ctx, PostCreate, req)
			return rs, err
		case NewMediaObject:
			s.MediaObject(ctx, req)
		case GetPost:
			{
				postId, _ := strconv.Atoi(req.Params.Param[0].Value.String)
				return s.GetPost(ctx, int64(postId))
			}
		}

		return resp, err
	}
}
