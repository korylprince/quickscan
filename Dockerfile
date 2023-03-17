# FROM golang:1-alpine as builder

# RUN go install github.com/korylprince/fileenv@v1.1.0

FROM alpine:latest

ARG GO_PROJECT_NAME
ENV GO_PROJECT_NAME=${GO_PROJECT_NAME}

RUN apk add --no-cache ca-certificates

COPY docker-entrypoint.sh /
COPY ${GO_PROJECT_NAME} /

CMD ["/docker-entrypoint.sh"]
