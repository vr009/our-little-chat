package delivery

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
	"our-little-chatik/internal/mailer/internal"
	"our-little-chatik/internal/models"
)

type RabbitListener struct {
	ch *amqp.Channel
	u  internal.MailerUsecase
}

func NewRabbitListener(ch *amqp.Channel,
	u internal.MailerUsecase) *RabbitListener {
	return &RabbitListener{
		ch: ch,
		u:  u,
	}
}

func (l RabbitListener) ListenToQueue() {
	q, err := l.ch.QueueDeclare(
		"mail", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)

	if err != nil {
		panic(err.Error())
	}

	msgs, err := l.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err.Error())
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			slog.Info("Received a message: %s", d.Body)
			task := models.ActivationTask{}
			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			err = l.u.SendMailMessage(task.Receiver, task)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
		}
	}()

	slog.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
