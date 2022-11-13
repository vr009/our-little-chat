FROM golang:latest

RUN go version

ENV GO111MODULE=on

ENV GOPATH=/

COPY . .

EXPOSE 8083

RUN go mod tidy

RUN go build -o chat-service ./internal/chat/cmd/main.go

CMD ["./chat-service"]
