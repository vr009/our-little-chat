FROM golang:latest

RUN go version

ENV GOPATH=/

COPY . .

EXPOSE 8086

RUN go mod tidy

RUN go build -o user-data-service ./internal/users/cmd/main.go

CMD ["./user-data-service"]
