FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

EXPOSE 8081

RUN apt update && apt install libssl-dev -y

RUN go mod download

RUN go mod tidy

RUN go build -o chat-diff-service ./internal/chat_diff/cmd/main.go

CMD ["./chat-diff-service"]
