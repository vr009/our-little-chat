FROM golang:latest

RUN go version

ENV GOPATH=/

COPY go.mod .

EXPOSE 8083

COPY . .

RUN go mod tidy

RUN go build -o chat-service ./internal/chat/cmd/main.go

CMD ["./chat-service"]
