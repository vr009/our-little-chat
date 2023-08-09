FROM golang:latest

RUN go version

ENV GO111MODULE=on

ENV GOPATH=/

COPY go.mod .

EXPOSE 8083

RUN go mod tidy

COPY . .

RUN go build -o chat-service ./internal/chat/cmd/main.go

CMD ["./chat-service"]
