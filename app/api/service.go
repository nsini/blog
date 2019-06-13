package api

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/config"
	"github.com/nsini/blog/repository"
	"strconv"
)

type Service interface {
	Post(ctx context.Context, method PostMethod, req postRequest) (rs newPostResponse, err error)
	GetPost(ctx context.Context, id int64) (rs map[string]interface{}, err error)
	GetCategories(ctx context.Context, req postRequest) (rs methodResponse, err error) // todo 需要调整 不应该让service返回xml
}

type service struct {
	post   repository.PostRepository
	user   repository.UserRepository
	image  repository.ImageRepository
	logger log.Logger
	config config.Config
}

func (c *service) Post(ctx context.Context, method PostMethod, req postRequest) (rs newPostResponse, err error) {

	_ = c.logger.Log("methodName", req.MethodName, "PostMethod", method, "username", req.Params.Param[1].Value.String, "password", req.Params.Param[2].Value.String)

	//c.post.Create(repository.Post{
	//
	//})

	for _, member := range req.Params.Param[3].Value.Struct.Member {
		_ = c.logger.Log("member", member.Name)
	}

	rs.Params.Param.Value.String = strconv.Itoa(20051)

	return
}

func (c *service) GetPost(ctx context.Context, id int64) (rs map[string]interface{}, err error) {

	_ = c.logger.Log("postId", id)

	return
}

func (c *service) GetCategories(ctx context.Context, req postRequest) (rs methodResponse, err error) {

	_ = c.logger.Log("methodName", req.MethodName)

	var data []dataValue

	var members, members0, members1 []member
	members = append(members, member{
		Name: "categoryId",
		Value: value{
			String: "1",
		},
	}, member{
		Name: "parentId",
		Value: value{
			String: "0",
		},
	}, member{
		Name: "categoryName",
		Value: value{
			String: "技术",
		},
	}, member{
		Name: "description",
		Value: value{
			String: "技术类的文章",
		},
	}, member{
		Name: "httpUrl",
		Value: value{
			String: "",
		},
	}, member{
		Name: "title",
		Value: value{
			String: "技术",
		},
	})

	members0 = append(members0, member{
		Name: "categoryId",
		Value: value{
			String: "2",
		},
	}, member{
		Name: "parentId",
		Value: value{
			String: "0",
		},
	}, member{
		Name: "categoryName",
		Value: value{
			String: "生活",
		},
	}, member{
		Name: "description",
		Value: value{
			String: "生活的文章",
		},
	}, member{
		Name: "httpUrl",
		Value: value{
			String: "",
		},
	}, member{
		Name: "title",
		Value: value{
			String: "生活",
		},
	})

	members1 = append(members1, member{
		Name: "categoryId",
		Value: value{
			String: "3",
		},
	}, member{
		Name: "parentId",
		Value: value{
			String: "0",
		},
	}, member{
		Name: "categoryName",
		Value: value{
			String: "旅游",
		},
	}, member{
		Name: "description",
		Value: value{
			String: "旅游的文章",
		},
	}, member{
		Name: "httpUrl",
		Value: value{
			String: "",
		},
	}, member{
		Name: "title",
		Value: value{
			String: "旅游",
		},
	})

	data = append(data, dataValue{
		Struct: valStruct{
			Member: members,
		},
	}, dataValue{
		Struct: valStruct{
			Member: members0,
		},
	}, dataValue{
		Struct: valStruct{
			Member: members1,
		},
	})

	rs.Params.Param.Value.Array.Data.Value = data

	return
}

func NewService(logger log.Logger, cf config.Config, post repository.PostRepository, user repository.UserRepository, image repository.ImageRepository) Service {
	return &service{
		post:   post,
		user:   user,
		image:  image,
		logger: logger,
		config: cf,
	}
}
