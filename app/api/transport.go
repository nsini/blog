package api

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/blog/repository"
	"github.com/nsini/blog/templates"
	"io/ioutil"
	"net/http"
)

var errBadRoute = errors.New("bad route")
var ErrInvalidArgument = errors.New("invalid argument")

func MakeHandler(ps Service, logger kitlog.Logger) http.Handler {
	//ctx := context.Background()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	post := kithttp.NewServer(
		makePostEndpoint(ps),
		decodePostRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/api/post/release", post).Methods("POST")
	return r
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req postRequest

	if err = xml.Unmarshal(b, &req); err != nil {
		// todo err 加一层处理
		return nil, err
	}

	bb, _ := xml.Marshal(req)

	fmt.Println(string(bb))

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	// todo 返回xml数据

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
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
