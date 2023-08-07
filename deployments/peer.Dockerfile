FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

EXPOSE 8084

#RUN go mod download

RUN go mod tidy

RUN go build -o peer-service ./internal/peer/cmd/main.go

CMD ["./peer-service"]
