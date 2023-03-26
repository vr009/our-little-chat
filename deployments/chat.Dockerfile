FROM golang:latest

RUN go version

ENV GO111MODULE=on

ENV GOPATH=/

COPY . .

EXPOSE 8083

RUN apt update && apt install libssl-dev -y

#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz |  \
#    tar xvz && mv migrate.linux-amd64 $GOPATH/bin/migrate
#
#RUN migrate -source file://internal/chat/db/migrations -database 'postgresql://localhost:5433/chats?user=service&password=test' up
#
#RUN migrate -source file://internal/chat/db/migrations/test -database 'postgresql://db-chat:5432/chats?user=service&password=test' up

RUN go mod tidy

RUN go build -o chat-service ./internal/chat/cmd/main.go

CMD ["./chat-service"]
