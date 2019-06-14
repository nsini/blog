package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/go-kit/kit/endpoint"
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

type newMediaObject struct {
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
							Text    string `xml:",chardata"`
							Boolean string `xml:"boolean"`
							Base64  string `xml:"base64"`
							String  string `xml:"string"`
						} `xml:"value"`
					} `xml:"member"`
				} `xml:"struct"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type methodResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				String string `xml:"string"`
				Text   string `xml:",chardata"`
				Array  struct {
					Text string `xml:",chardata"`
					Data struct {
						Text  string      `xml:",chardata"`
						Value []dataValue `xml:"value"`
					} `xml:"data"`
				} `xml:"array"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type member struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name"`
	Value value  `xml:"value"`
}

type value struct {
	Text   string `xml:",chardata"`
	String string `xml:"string"`
}

type valStruct struct {
	Text   string   `xml:",chardata"`
	Member []member `xml:"member"`
}

type dataValue struct {
	Text   string    `xml:",chardata"`
	Struct valStruct `xml:"struct"`
}

type newPostResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				String string `xml:"string"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		resp := methodResponse{}
		var err error

		// todo 进行token验证
		fmt.Println("username", req.Params.Param[1].Value.String)
		fmt.Println("password", req.Params.Param[2].Value.String)

		switch PostMethod(req.MethodName) {
		case GetUsersBlogs:
			// todo check response
			//resp.Params.Param.Value.Array.Data.Value.Struct.Member = append(resp.Params.Param.Value.Array.Data.Value.Struct.Member, member{
			//	Name: "blogid",
			//	Value: value{
			//		String: "dudulu",
			//	},
			//})
			//resp.Params.Param.Value.Array.Data.Value.Struct.Member = append(resp.Params.Param.Value.Array.Data.Value.Struct.Member, member{
			//	Name: "blogName",
			//	Value: value{
			//		String: "nsini",
			//	},
			//})
		case GetCategories:
			return s.GetCategories(ctx, req)
		case PostCreate:
			rs, err := s.Post(ctx, PostCreate, req)
			return rs, err
		case GetPost:
			{
				postId, _ := strconv.Atoi(req.Params.Param[0].Value.String)
				s.GetPost(ctx, int64(postId))
			}
		}

		return resp, err
	}
}
