package api

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/config"
	"github.com/nsini/blog/repository"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
	"strconv"
	"strings"
	"time"
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

type PostFields string

const (
	PostStatus      PostFields = "post_status"
	PostType        PostFields = "post_type"
	PostCategories  PostFields = "categories"
	PostTitle       PostFields = "title"
	PostDateCreated PostFields = "dateCreated"
	PostWpSlug      PostFields = "wp_slug"
	PostDescription PostFields = "description"
	PostKeywords    PostFields = "mt_keywords"
)

func (c *service) Post(ctx context.Context, method PostMethod, req postRequest) (rs newPostResponse, err error) {

	_ = c.logger.Log("methodName", req.MethodName, "PostMethod", method, "username", req.Params.Param[1].Value.String, "password", req.Params.Param[2].Value.String)

	//c.post.Create(repository.Post{
	//
	//})

	var isMarkdown bool
	var postStatus, postType, postTitle, slug, description string
	var categories []string
	var keywords []string
	var postDateCreated time.Time

	for _, member := range req.Params.Param[3].Value.Struct.Member {
		_ = c.logger.Log("member", member.Name)
		switch PostFields(member.Name) {
		case PostStatus:
			postStatus = member.Value.String
		case PostType:
			postType = member.Value.String
		case PostCategories:
			for _, val := range member.Value.Array.Data {
				categories = append(categories, val.Value.String)
			}
		case PostTitle:
			postTitle = member.Value.String
		case PostDateCreated:
			load, _ := time.LoadLocation("Asia/Shanghai")
			if postDateCreated, err = time.ParseInLocation("20060102T15:04:05Z", member.Value.DateTimeIso8601, load); err == nil {
				_ = c.logger.Log("time", "Parse", "err", err)
				postDateCreated = postDateCreated.Add(8 * 3600 * time.Second)
			} else {
				postDateCreated = time.Now()
			}
		case PostWpSlug:
			slug = member.Value.String
		case PostDescription:
			description = member.Value.String
		case PostKeywords:
			keywords = strings.Split(member.Value.String, ",")
		}
	}

	_ = c.logger.Log("req.Params.Param[4].Value.Boolean", req.Params.Param[4].Value.Boolean)

	if boolean, _ := strconv.Atoi(req.Params.Param[4].Value.Boolean); boolean == 1 {
		fmt.Println(boolean)
		isMarkdown = true
	}

	markdown, _ := strconv.Atoi(req.Params.Param[4].Value.Boolean)

	_ = c.logger.Log("postStatus", postStatus, "postType", postType, "categories", categories, "postDateCreated", postDateCreated.Format("2006-01-02 15:04:05"), "postTitle", postTitle, "slug", slug, "description", description, "keywords", keywords)

	_ = c.logger.Log("isMarkdown", isMarkdown)

	// todo 查询用户获取用户ID
	userId := int64(1)

	if err = c.post.Create(repository.Post{
		Title:       postTitle,
		Content:     description,
		Description: null.StringFrom(description[:100]),
		IsMarkdown:  null.IntFrom(int64(markdown)),
		PushTime:    null.NewTime(time.Now(), true),
		UserID:      null.IntFrom(userId),
		Status:      1,
		Action:      1,
	}); err != nil {
		return
	}

	return rs, errors.New("test")

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
