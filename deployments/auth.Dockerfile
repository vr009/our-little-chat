FROM golang:latest

RUN go version

ENV GO111MODULE=on

ENV GOPATH=/

COPY . .

EXPOSE 8087

RUN go mod tidy

RUN go build -o auth-service ./internal/auth/cmd/main.go


CMD ["./auth-service"]
