FROM golang:1.19 AS build
WORKDIR $GOPATH/src/github.com/codeinuit/shibesbot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/shibesbot

FROM alpine:latest
WORKDIR /root/
COPY --from=build /go/src/github.com/codeinuit/shibesbot/app .
ENTRYPOINT ["./app"]
