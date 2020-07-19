FROM golang:1.14.4 AS build
WORKDIR $GOPATH/src/github.com/P147x/shibesbot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=build /go/src/github.com/P147x/shibesbot/app .
ENTRYPOINT ["./app"]
