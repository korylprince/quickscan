FROM golang:1-alpine as builder

ARG VERSION

RUN apk add --no-cache sqlite build-base

RUN go install "github.com/korylprince/quickscan@$VERSION"

FROM alpine:3

RUN apk add --no-cache ca-certificates sqlite

COPY --from=builder /go/bin/quickscan /

CMD ["/quickscan"]
