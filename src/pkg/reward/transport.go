/**
 * @Time: 2019-07-20 16:38
 * @Author: solacowa@gmail.com
 * @File: transport
 * @Software: GoLand
 */

package reward

import (
	"context"
	"errors"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/blog/src/repository"
	"github.com/nsini/blog/src/templates"
	"net/http"
)

var errBadRoute = errors.New("bad route")

func MakeHandler(logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := mux.NewRouter()
	r.Handle("/reward", kithttp.NewServer(
		func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return nil, nil
		},
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
			if e, ok := response.(errorer); ok && e.error() != nil {
				encodeError(ctx, e.error(), w)
				return nil
			}

			ctx = context.WithValue(ctx, "method", "reward")
			return templates.RenderHtml(ctx, w, nil)
		},
		opts...,
	)).Methods("GET")

	return r
}

type errorer interface {
	error() error
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case repository.PostNotFound:
		w.WriteHeader(http.StatusNotFound)
		ctx = context.WithValue(ctx, "method", "404")
		_ = templates.RenderHtml(ctx, w, map[string]interface{}{})
		return
	case errBadRoute:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, _ = w.Write([]byte(err.Error()))
}
