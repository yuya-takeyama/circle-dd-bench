FROM golang:1.11-alpine3.8 as builder

RUN apk --update add gcc

WORKDIR /go/src/github.com/yuya-takeyama/circle-dd-bench
COPY . /go/src/github.com/yuya-takeyama/circle-dd-bench

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"'

FROM alpine:3.10.3

COPY --from=builder /go/src/github.com/yuya-takeyama/circle-dd-bench/circle-dd-bench /usr/local/bin

ENTRYPOINT ["circle-dd-bench"]
