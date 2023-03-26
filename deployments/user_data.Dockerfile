FROM golang:latest

RUN go version

ENV GO111MODULE=on

ENV GOPATH=/

COPY . .

EXPOSE 8086

RUN go mod tidy

RUN apt-get update && apt-get install libssl-dev -y

RUN go build -o user-data-service ./internal/user_data/cmd/main.go

CMD ["./user-data-service"]
