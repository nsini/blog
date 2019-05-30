package post

import (
	"context"
	"errors"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/blog/repository"
	"github.com/nsini/blog/templates"
	"net/http"
	"strconv"
)

var errBadRoute = errors.New("bad route")

func MakeHandler(ps Service, logger kitlog.Logger) http.Handler {
	//ctx := context.Background()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	detail := kithttp.NewServer(
		makeDetailEndpoint(ps),
		decodeDetailRequest,
		encodeDetailResponse,
		opts...,
	)

	list := kithttp.NewServer(
		makeListEndpoint(ps),
		decodeListRequest,
		encodeListResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/post/", list).Methods("GET")
	r.Handle("/post/{id}", detail).Methods("GET")
	return r
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	size, ok := vars["pageSize"]
	if !ok {
		size = "10"
	}
	order, ok := vars["order"]
	if !ok {
		order = "desc"
	}
	by, ok := vars["by"]
	if !ok {
		by = "id"
	}
	limit, ok := vars["limit"]
	if !ok {
		limit = "10"
	}
	offset, ok := vars["offset"]
	if !ok {
		offset = "0"
	}

	pageSize, _ := strconv.Atoi(size)
	pageLimit, _ := strconv.Atoi(limit)
	pageOffset, _ := strconv.Atoi(offset)
	return listRequest{
		pageSize: pageSize,
		order:    order,
		by:       by,
		limit:    pageLimit,
		offset:   pageOffset,
	}, nil
}

func decodeDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	postId, err := strconv.Atoi(id)
	if err != nil {
		return nil, errBadRoute
	}
	return postRequest{
		Id: int64(postId),
	}, nil
}

func encodeDetailResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "blog-single")

	resp := response.(postResponse)

	return templates.RenderHtml(ctx, w, map[string]interface{}{
		"content":    resp.Data.Content,
		"title":      resp.Data.Title,
		"publish_at": resp.Data.PushTime.Time.Format("2006/01/02 15:04:05"),
		"updated_at": resp.Data.UpdatedAt,
		"author":     "嘟嘟噜",
		"comment":    4,
	})
}

func encodeListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "blog-left-sidebar")

	resp := response.(listResponse)

	fmt.Println(resp)
	return templates.RenderHtml(ctx, w, map[string]interface{}{
		"list": resp.Data,
	})
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
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, _ = w.Write([]byte(err.Error()))
}
