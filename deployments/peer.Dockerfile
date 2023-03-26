FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

EXPOSE 8082

RUN go mod download

RUN go mod tidy

RUN apt-get update && apt-get install libssl-dev -y

RUN go build -o peer-service ./internal/peer/cmd/main.go

CMD ["./peer-service"]
