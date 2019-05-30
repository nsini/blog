FROM golang:latest as build-dev

ENV GO111MODULE=on
ENV BUILDPATH=github.com/nsini/blog
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
RUN cd /go/src/${BUILDPATH} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v

FROM alpine:latest

COPY --from=build-env /go/bin/blog /go/bin/blog
WORKDIR /go/bin/
CMD ["/go/bin/blog"]