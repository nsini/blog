package api

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/blog/src/repository"
	"io/ioutil"
	"net/http"
)

//var errBadRoute = errors.New("bad route")
var ErrInvalidArgument = errors.New("invalid argument")

func MakeHandler(ps Service, logger kitlog.Logger) http.Handler {
	//ctx := context.Background()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeXmlError),
	}

	post := kithttp.NewServer(
		makePostEndpoint(ps),
		decodePostRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/api/post/metaweblog", post).Methods("POST")
	return r
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))

	var req postRequest

	if err = xml.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	switch req.MethodName {
	case NewMediaObject.String():
		{
			var req postRequest
			if err = xml.Unmarshal(b, &req); err != nil {
				return nil, err
			}
			return req, nil
		}
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeXmlError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	return xml.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeXmlError(ctx context.Context, err error, w http.ResponseWriter) {

	type faultStruct struct {
		Text   string    `xml:",chardata"`
		Struct valStruct `xml:"struct"`
	}

	type fault struct {
		Text  string      `xml:",chardata"`
		Value faultStruct `xml:"value"`
	}

	type errorResponse struct {
		XMLName xml.Name `xml:"methodResponse"`
		Text    string   `xml:",chardata"`
		Fault   fault    `xml:"fault"`
	}

	var faultCode string

	switch err {
	case repository.PostNotFound:
		faultCode = "404"
	case ErrInvalidArgument:
		faultCode = "401"
	case NoPermission:
		faultCode = "403"
	default:
		faultCode = "500"
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_ = xml.NewEncoder(w).Encode(errorResponse{
		Fault: fault{
			Value: faultStruct{
				Struct: valStruct{
					Member: []member{
						{Name: "faultString", Value: memberValue{String: faultCode}},
						{Name: "faultCode", Value: memberValue{String: err.Error()}},
					},
				},
			},
		},
	})
}
