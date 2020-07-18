FROM golang:1.14.4 AS build
COPY . /go/src/github.com/P147x/shibesbot
WORKDIR /go/src/github.com/P147x/shibesbot
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./src/

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/P147x/shibesbot/app .
ENTRYPOINT ["./app"]
