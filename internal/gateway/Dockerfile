FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

EXPOSE 8000

RUN go mod download

RUN go mod tidy

RUN go build -o chat_history-service ./cmd/main.go


CMD ["./chat-service"]
