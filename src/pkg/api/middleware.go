/**
 * @Time : 2019-09-10 10:32
 * @Author : solacowa@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package api

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/blog/src/repository"
	"github.com/nsini/blog/src/util/encode"
)

type ASDContext string

const (
	UserIdContext ASDContext = "user-id"
)

var ErrorASD = errors.New("权限验证失败！")

func checkAuthMiddleware(logger log.Logger, repository repository.Repository, salt string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(postRequest)
			username := req.Params.Param[1].Value.String
			password := req.Params.Param[2].Value.String

			_ = level.Debug(logger).Log("username", username, "password", password)

			if username == "" || password == "" {
				_ = level.Error(logger).Log("username or password:", "is nil")
				return nil, ErrorASD
			}

			user, err := repository.User().FindAndPwd(username, encode.EncodePassword(password, salt))
			if err != nil {
				_ = level.Error(logger).Log("User", "FindAndPwd", "err", err.Error())
				return nil, ErrorASD
			}

			ctx = context.WithValue(ctx, UserIdContext, user.ID)

			return next(ctx, request)
		}
	}
}

func imageCheckAuthMiddleware(logger log.Logger, repository repository.Repository, salt string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(uploadImageRequest)
			if req.Username == "" || req.Password == "" {
				_ = level.Error(logger).Log("username or password:", "is nil")
				return nil, ErrorASD
			}

			_ = level.Debug(logger).Log("username", req.Username, "password", req.Password)

			user, err := repository.User().FindAndPwd(req.Username, encode.EncodePassword(req.Password, salt))
			if err != nil {
				_ = level.Error(logger).Log("User", "FindAndPwd", "err", err.Error())
				return nil, ErrorASD
			}

			ctx = context.WithValue(ctx, UserIdContext, user.ID)

			return next(ctx, request)
		}
	}
}
