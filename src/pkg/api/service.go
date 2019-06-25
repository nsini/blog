package api

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/src/config"
	"github.com/nsini/blog/src/repository"
	"github.com/nsini/blog/src/tools"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"strconv"
	"strings"
	"time"
)

var PostNotFound = errors.New("post not found!")

type Service interface {
	Authentication(ctx context.Context, req postRequest) (rs getUsersBlogsResponse, err error)
	Post(ctx context.Context, method PostMethod, req postRequest) (rs newPostResponse, err error)
	GetPost(ctx context.Context, id int64) (rs *getPostResponse, err error)
	GetCategories(ctx context.Context, req postRequest) (rs *getCategoriesResponse, err error) // todo 需要调整 不应该让service返回xml
	MediaObject(ctx context.Context, req postRequest)
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
	MediaOverwrite  PostFields = "overwrite"
	MediaBits       PostFields = "bits"
	MediaName       PostFields = "name"
	MediaType       PostFields = "type"
)

func (c *service) Authentication(ctx context.Context, req postRequest) (rs getUsersBlogsResponse, err error) {

	return
}

/**
 * @Title 上传流媒体类型的文件
 */
func (c *service) MediaObject(ctx context.Context, req postRequest) {
	var overwrite bool
	var bits, mediaName, mediaType string
	for _, val := range req.Params.Param[3].Value.Struct.Member {
		_ = c.logger.Log("val", val.Name)
		switch PostFields(val.Name) {
		case MediaOverwrite:
			overwrite, _ = strconv.ParseBool(val.Value.Boolean)
		case MediaBits:
			bits = val.Value.Base64
		case MediaName:
			mediaName = strings.TrimSpace(strings.ToLower(val.Value.String))
		case MediaType:
			mediaType = val.Value.String
		}
	}

	bits = strings.TrimSpace(strings.Trim(bits, "\n"))
	bits = strings.Replace(bits, " ", "", -1)
	dist, err := base64.StdEncoding.DecodeString(bits)

	if err != nil {
		_ = c.logger.Log("base64", "DecodeString", "err", err.Error())
		return
	}

	if err = ioutil.WriteFile("/tmp/"+mediaName, dist, 0666); err != nil {
		_ = c.logger.Log("ioutil", "WriteFile", "err", err.Error())
		return
	}

	// 先存，再验证，失败再删除
	f, err := os.Open("/tmp/" + mediaName)
	if err != nil {
		_ = c.logger.Log("os", "Open", "err", err.Error())
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			_ = c.logger.Log("f", "Close", "err", err.Error())
		}
	}()

	md5h := md5.New()
	if _, err = io.Copy(md5h, f); err != nil {
		_ = c.logger.Log("io", "Copy", "err", err.Error())
		return
	}

	var fileSize int64
	if fileInfo, err := os.Stat("/tmp/" + mediaName); err == nil {
		fileSize = fileInfo.Size()
	}

	//defer func() {
	//	if err = os.Remove("/tmp/"+mediaName); err != nil {
	//		_ = c.logger.Log("os", "Remove", "err", err.Error())
	//	}
	//}()

	fileSha := fmt.Sprintf("%x", md5h.Sum([]byte("")))

	// 进行数据md5值验证 需要不需要返回地址呢？
	if c.image.ExistsImageByMd5(fileSha) {
		_ = c.logger.Log("c.image", "ExistsImageByMd5", "err", "file is exists.")
		return
	}

	simPath := time.Now().Format("2006/01/") + fileSha[len(fileSha)-5:len(fileSha)-3] + "/" + fileSha[24:26] + "/" + fileSha[16:17] + fileSha[12:13] + "/"
	filePath := c.config.Get(config.ImageFilePath) + simPath
	if !tools.PathExist(filePath) {
		if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
			_ = c.logger.Log("os", "MkdirAll", "err", err.Error())
			return
		}
	}

	var extName = ".jpg"
	if exts, err := mime.ExtensionsByType(mediaType); err == nil {
		extName = exts[0]
	}

	fileName := time.Now().Format("20060102") + "-" + fileSha + extName
	fileFullPath := filePath + fileName

	if err = os.Rename("/tmp/"+mediaName, fileFullPath); err != nil {
		_ = c.logger.Log("os", "Rename", "err", err.Error())
		return
	}

	// 存入数据库
	if err = c.image.AddImage(&repository.Image{
		ImageName: fileName,
		Extension: null.StringFrom(extName),
		ImagePath: null.StringFrom(simPath + fileName),
		RealPath:  null.StringFrom(fileFullPath),
		//ImageTime:          null.NewTime(time.Now(), false),
		ImageStatus:        null.IntFrom(0),
		ImageSize:          null.StringFrom(strconv.Itoa(int(fileSize))),
		Md5:                null.StringFrom(fileSha),
		ClientOriginalMame: null.StringFrom(mediaName),
	}); err != nil {
		_ = c.logger.Log("c.image", "AddImage", "err", err.Error())
		return
	}

	// todo 返回图片的xml response

	_ = c.logger.Log("overwrite", overwrite, "bits", "", "mediaName", mediaName, "mediaType", mediaType, "fileSha", fileSha, "fileName", fileName)
}

/**
 * @Title 发布内容
 */
func (c *service) Post(ctx context.Context, method PostMethod, req postRequest) (rs newPostResponse, err error) {

	_ = c.logger.Log("methodName", req.MethodName, "PostMethod", method, "username", req.Params.Param[1].Value.String, "password", req.Params.Param[2].Value.String)

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

	publishStatus, _ := strconv.Atoi(req.Params.Param[4].Value.Boolean) // todo 1: 已发布，0: 草稿

	_ = c.logger.Log("postStatus", postStatus, "postType", postType, "categories", categories, "postDateCreated", postDateCreated.Format("2006-01-02 15:04:05"), "postTitle", postTitle, "slug", slug, "keywords", keywords)

	// todo 查询用户获取用户ID
	userId := int64(1)

	desc := []rune(description)
	if len(desc) > 100 {
		desc = desc[:100]
	}

	p := repository.Post{
		Title:       postTitle,
		Content:     description,
		Description: null.StringFrom(string(desc)),
		IsMarkdown:  null.IntFrom(int64(1)), // todo 想办法怎么验证一下
		PushTime:    null.NewTime(time.Now(), true),
		UserID:      null.IntFrom(userId),
		Status:      publishStatus,
		Action:      1,
		ReadNum:     1,
	}

	if err = c.post.Create(&p); err != nil {
		return
	}

	rs.Params.Param.Value.String = strconv.Itoa(int(p.Model.ID))

	return
}

/**
 * @Title 获取文章
 */
func (c *service) GetPost(ctx context.Context, id int64) (rs *getPostResponse, err error) {

	_ = c.logger.Log("postId", id)

	post, err := c.post.Find(id)
	if err != nil {
		return nil, PostNotFound
	}

	var members []member

	members = append(members, member{
		Name: "userid",
		Value: memberValue{
			String: strconv.Itoa(int(post.UserID.Int64)),
		},
	}, member{
		Name: "postid",
		Value: memberValue{
			String: strconv.Itoa(int(post.Model.ID)),
		},
	}, member{
		Name: "description",
		Value: memberValue{
			String: post.Description.String,
		},
	}, member{
		Name: "title",
		Value: memberValue{
			String: post.Title,
		},
	}, member{
		Name: "link",
		Value: memberValue{
			String: "/post/" + strconv.Itoa(int(post.Model.ID)),
		},
	}, member{
		Name: "mt_keywords",
		Value: memberValue{
			String: "存储,Golang",
		},
	}, member{
		Name: "wp_slug",
		Value: memberValue{
			String: post.Slug.String,
		},
	}, member{
		Name: "wp_author",
		Value: memberValue{
			String: "",
		},
	}, member{
		Name: "wp_author_id",
		Value: memberValue{
			String: "",
		},
	}, member{
		Name: "date_created_gmt",
		Value: memberValue{
			String: post.PushTime.Time.String(),
		},
	}, member{
		Name: "post_status",
		Value: memberValue{
			String: strconv.Itoa(post.Status),
		},
	}, member{
		Name: "categories",
		//Value: memberValue{
		//	Array: array{
		//		Data: data{
		//			Value: dataValue{
		//				String: "技术,生活",
		//			},
		//		},
		//	},
		//},
	})

	resp := new(getPostResponse)

	resp.Params.Param.Value.Struct.Member = members

	return resp, nil
}

/**
 * @Title 获取分类列表
 */
func (c *service) GetCategories(ctx context.Context, req postRequest) (rs *getCategoriesResponse, err error) {

	_ = c.logger.Log("methodName", req.MethodName)

	resp := new(getCategoriesResponse)
	resp.Params.Param.Value.Array.Data.Value = append(resp.Params.Param.Value.Array.Data.Value, dataValue{
		Struct: valStruct{
			Member: []member{
				{Name: "categoryId", Value: memberValue{String: "1"}},
				{Name: "parentId", Value: memberValue{String: "0"}},
				{Name: "categoryName", Value: memberValue{String: "技术"}},
				{Name: "description", Value: memberValue{String: "技术类的文章"}},
				{Name: "title", Value: memberValue{String: "技术文章"}},
			},
		},
	}, dataValue{
		Struct: valStruct{
			Member: []member{
				{Name: "categoryId", Value: memberValue{String: "2"}},
				{Name: "parentId", Value: memberValue{String: "0"}},
				{Name: "categoryName", Value: memberValue{String: "生活"}},
				{Name: "description", Value: memberValue{String: "生活类的文章"}},
				{Name: "title", Value: memberValue{String: "生活文章"}},
			},
		},
	}, dataValue{
		Struct: valStruct{
			Member: []member{
				{Name: "categoryId", Value: memberValue{String: "3"}},
				{Name: "parentId", Value: memberValue{String: "0"}},
				{Name: "categoryName", Value: memberValue{String: "旅游"}},
				{Name: "description", Value: memberValue{String: "旅游的文章"}},
				{Name: "title", Value: memberValue{String: "旅游文章"}},
			},
		},
	})

	return resp, nil
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
