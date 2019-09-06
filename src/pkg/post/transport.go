package post

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/blog/src/repository"
	"github.com/nsini/blog/src/templates"
	"net/http"
	"strconv"
	"strings"
)

var errBadRoute = errors.New("bad route")

func MakeHandler(ps Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	detail := kithttp.NewServer(
		makeGetEndpoint(ps),
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

	popular := kithttp.NewServer(
		makePopularEndpoint(ps),
		decodePopularRequest,
		encodePopularResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/post", list).Methods("GET")
	r.Handle("/post/popular", popular).Methods("GET")
	r.Handle("/post/{id:[0-9]+}", detail).Methods("GET")
	return r
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
		by = "id"
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

	ctx = context.WithValue(ctx, "method", "blog-single")

	resp := response.(postResponse)

	return templates.RenderHtml(ctx, w, resp.Data)
}

func encodeListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "blog-left-sidebar")

	resp := response.(listResponse)

	return templates.RenderHtml(ctx, w, map[string]interface{}{
		"list":      resp.Data,
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
	res = append(res, fmt.Sprintf(`<li><a href="/post?pageSize=10&offset=%d">Prev</a></li>`, prev))
	for i := 1; i < (count / pageSize); i++ {
		os := (i - 1) * 10
		var active string
		if offset == os {
			active = `class="active"`
		}
		res = append(res, fmt.Sprintf(`<li %s><a href="/post?pageSize=10&offset=%d">%d</a></li>`, active, os, i))
	}
	res = append(res, fmt.Sprintf(`<li><a href="/post?pageSize=10&offset=%d">Next</a></li>`, next))
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