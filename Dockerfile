FROM golang:1.13.0-alpine3.10 as build-env

ENV GO111MODULE=on
ENV GOPROXY=http://goproxy.yrd.creditease.corp
ENV BUILDPATH=github.com/nsini/blog
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
RUN cd /go/src/${BUILDPATH} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v

FROM alpine:latest

COPY --from=build-env /go/bin/blog /go/bin/blog
COPY ./views /go/bin/
COPY ./static /go/bin/
WORKDIR /go/bin/
CMD ["/go/bin/blog"]
