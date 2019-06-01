package post

import (
	"context"
	"encoding/json"
	"errors"
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

	popular := kithttp.NewServer(
		makePopularEndpoint(ps),
		decodePopularRequest,
		encodePopularResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/post/", list).Methods("GET")
	r.Handle("/post/popular", popular).Methods("GET")
	r.Handle("/post/{id}", detail).Methods("GET")
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
	//ctx = context.WithValue(ctx, "pageSize", pageSize)
	//ctx = context.WithValue(ctx, "limit", pageLimit)
	//ctx = context.WithValue(ctx, "offset", pageOffset)
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
		"list":  resp.Data,
		"count": resp.Count,
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
