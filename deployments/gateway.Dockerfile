FROM golang:latest

RUN go version

ENV GO111MODULE=on

ENV GOPATH=/

COPY . .

EXPOSE 8080

RUN go mod tidy

RUN go build -o gateway-service ./internal/gateway/cmd/main.go

CMD ["./gateway-service"]
