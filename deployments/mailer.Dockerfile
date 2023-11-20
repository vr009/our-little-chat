FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

EXPOSE 8088

#RUN go mod download

RUN go mod tidy

RUN go build -o flusher-service ./internal/mailer/cmd/main.go

CMD ["./mailer-service"]
