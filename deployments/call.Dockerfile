FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

EXPOSE 8090

#RUN go mod download

RUN go mod tidy

RUN go build -o call-service ./internal/call/cmd/main.go

CMD ["./call-service"]
