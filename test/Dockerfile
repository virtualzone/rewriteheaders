FROM golang:1.17-alpine
WORKDIR /go/src/app
ADD testserver.go .
ADD go.mod .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .
EXPOSE 8081
CMD ["./main"]