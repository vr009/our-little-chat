package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
	"os"
	"our-little-chatik/internal/mailer/internal/delivery"
	"our-little-chatik/internal/mailer/internal/repo"
	"our-little-chatik/internal/mailer/internal/usecase"
	"strconv"
)

type queueOpts struct {
	url string
}

type smtpOpts struct {
	host             string
	username         string
	password         string
	sender           string
	templateFilePath string
	port             int
}

type serverOpts struct {
	qOpts    queueOpts
	smtpOpts smtpOpts
}

func main() {
	opts, err := lookUpOpts()
	if err != nil {
		panic(err)
	}

	conn, err := amqp.Dial(opts.qOpts.url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	smtpConnector := repo.NewDefaultSMTPConnector(opts.smtpOpts.host,
		opts.smtpOpts.port,
		opts.smtpOpts.username,
		opts.smtpOpts.password,
		opts.smtpOpts.sender)
	smtpUsecase := usecase.NewDefaultMailUsecase(smtpConnector, opts.smtpOpts.templateFilePath)
	queue := delivery.NewRabbitListener(ch, smtpUsecase)

	slog.Info("Started server with opts: %v", opts)
	queue.ListenToQueue()
}

func lookUpOpts() (*serverOpts, error) {
	var err error
	errFmtMsg := "failed to get %s from env"
	opts := &serverOpts{}
	val := os.Getenv("QUEUE_URL")
	if val == "" {
		return nil, fmt.Errorf(errFmtMsg, "QUEUE_URL")
	}
	opts.qOpts.url = val

	val = os.Getenv("SMTP_HOST")
	if val == "" {
		return nil, fmt.Errorf(errFmtMsg, "SMTP_HOST")
	}
	opts.smtpOpts.host = val

	val = os.Getenv("SMTP_USERNAME")
	if val == "" {
		return nil, fmt.Errorf(errFmtMsg, "SMTP_USERNAME")
	}
	opts.smtpOpts.username = val

	val = os.Getenv("SMTP_PASSWORD")
	if val == "" {
		return nil, fmt.Errorf(errFmtMsg, "SMTP_PASSWORD")
	}
	opts.smtpOpts.password = val

	val = os.Getenv("SMTP_SENDER")
	if val == "" {
		return nil, fmt.Errorf(errFmtMsg, "SMTP_SENDER")
	}
	opts.smtpOpts.sender = val

	opts.smtpOpts.templateFilePath = "templates/user_welcome.tmpl"

	val = os.Getenv("SMTP_PORT")
	if val == "" {
		return nil, fmt.Errorf(errFmtMsg, "SMTP_PORT")
	}

	opts.smtpOpts.port, err = strconv.Atoi(val)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
