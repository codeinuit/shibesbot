FROM golang:1.16.5 AS build
WORKDIR $GOPATH/src/github.com/P147x/shibesbot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/shibesbot

FROM alpine:latest
WORKDIR /root/
COPY --from=build /go/src/github.com/P147x/shibesbot/app .
ENTRYPOINT ["./app"]
