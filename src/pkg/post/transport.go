package post

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/blog/src/encode"
	"github.com/nsini/blog/src/repository"
	"github.com/nsini/blog/src/templates"
	"golang.org/x/time/rate"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var errBadRoute = errors.New("bad route")

const rateBucketNum = 2

func MakeHandler(ps Service, logger kitlog.Logger, repository repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	ems := []endpoint.Middleware{
		UpdatePostReadNum(logger, repository),
		newTokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)),
	}

	mw := map[string][]endpoint.Middleware{
		"Get":     ems[:1],
		"Awesome": ems[1:],
	}

	eps := NewEndpoint(ps, mw)

	detail := kithttp.NewServer(
		eps.GetEndpoint,
		decodeDetailRequest,
		encodeDetailResponse,
		opts...,
	)

	list := kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encodeListResponse,
		opts...,
	)

	popular := kithttp.NewServer(
		eps.PopularEndpoint,
		decodePopularRequest,
		encodePopularResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/post", list).Methods(http.MethodGet)
	r.Handle("/post/popular", popular).Methods(http.MethodGet)
	r.Handle("/post/{id:[0-9]+}", detail).Methods(http.MethodGet)
	r.Handle("/post/{id:[0-9]+}",
		kithttp.NewServer(
			eps.AwesomeEndpoint,
			decodeAwesomeRequest,
			encode.EncodeJsonResponse,
			opts...,
		)).Methods(http.MethodPut)
	return r
}

func decodeAwesomeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	return postRequest{Id: postId}, nil
}

func decodePopularRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return popularRequest{}, nil
}

func encodePopularResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	size := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	by := r.URL.Query().Get("by")
	offset := r.URL.Query().Get("offset")
	action, _ := strconv.Atoi(r.URL.Query().Get("action"))

	if size == "" {
		size = "10"
	}
	if order == "" {
		order = "desc"
	}
	if by == "" {
		by = "push_time"
	}
	if offset == "" {
		offset = "0"
	}
	if action < 1 {
		action = 1
	}

	pageSize, _ := strconv.Atoi(size)
	pageOffset, _ := strconv.Atoi(offset)
	return listRequest{
		pageSize: pageSize,
		order:    order,
		by:       by,
		offset:   pageOffset,
		action:   action,
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

	ctx = context.WithValue(ctx, "method", "info")

	resp := response.(postResponse)

	return templates.RenderHtml(ctx, w, resp.Data)
}

func encodeListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "list")

	resp := response.(listResponse)

	other := resp.Data["other"].(map[string]interface{})

	return templates.RenderHtml(ctx, w, map[string]interface{}{
		"list":      resp.Data["post"],
		"tags":      other["tags"],
		"populars":  other["populars"],
		"total":     strconv.Itoa(int(resp.Count)),
		"paginator": postPaginator(int(resp.Count), resp.Paginator.PageSize, resp.Paginator.Offset),
	})
}

func postPaginator(count, pageSize, offset int) string {
	var res []string
	var prev, next int
	prev = offset - pageSize
	next = offset + pageSize
	if offset-pageSize < 0 {
		prev = 0
	}
	if offset+pageSize > count {
		next = offset
	}
	res = append(res, fmt.Sprintf(`<a href="/post?pageSize=10&offset=%d">上一页</a>&nbsp;`, prev))

	length := math.Ceil(float64(count) / float64(pageSize))
	for i := 1; i <= int(length); i++ {
		os := (i - 1) * 10
		if offset == os {
			res = append(res, fmt.Sprintf(`<b>%d</b>`, i))
			continue
		}
		res = append(res, fmt.Sprintf(`<a href="/post?pageSize=10&offset=%d">%d</a>&nbsp;`, os, i))
	}
	res = append(res, fmt.Sprintf(`<a href="/post?pageSize=10&offset=%d">下一页</a>&nbsp;`, next))
	return strings.Join(res, "\n")
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
