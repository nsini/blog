package api

import (
	"context"
	"encoding/xml"
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

type postRequest struct {
	XMLName    xml.Name `xml:"methodCall"`
	Text       string   `xml:",chardata"`
	MethodName string   `xml:"methodName"`
	Params     struct {
		Text  string `xml:",chardata"`
		Param []struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				String string `xml:"string"`
				Struct struct {
					Text   string `xml:",chardata"`
					Member []struct {
						Text  string `xml:",chardata"`
						Name  string `xml:"name"`
						Value struct {
							Text   string `xml:",chardata"`
							String string `xml:"string"`
							Array  struct {
								Text string `xml:",chardata"`
								Data []struct {
									Text  string `xml:",chardata"`
									Value struct {
										Text   string `xml:",chardata"`
										String string `xml:"string"`
									} `xml:"value"`
								} `xml:"data"`
							} `xml:"array"`
							DateTimeIso8601 string `xml:"dateTime.iso8601"`
						} `xml:"value"`
					} `xml:"member"`
				} `xml:"struct"`
				Boolean string `xml:"boolean"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
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
		case GetPost:
			{
				postId, _ := strconv.Atoi(req.Params.Param[0].Value.String)
				return s.GetPost(ctx, int64(postId))
			}
		}

		return resp, err
	}
}
