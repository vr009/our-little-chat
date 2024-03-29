FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

EXPOSE 8082

#RUN go mod download

RUN go mod tidy

RUN go build -o flusher-service ./internal/flusher/cmd/main.go

CMD ["./flusher-service"]
