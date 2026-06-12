FROM golang:1.21 AS build
WORKDIR $GOPATH/src/github.com/codeinuit/shibesbot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/shibesbot

FROM alpine:3.19
WORKDIR /root/
COPY --from=build /go/src/github.com/codeinuit/shibesbot/app .
ENTRYPOINT ["./app"]
