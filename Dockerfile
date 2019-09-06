FROM golang:latest as build-env

ENV GO111MODULE=on
ENV GOPROXY=http://goproxy-operations.kpl.yixinonline.org
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