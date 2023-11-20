package repo

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
	internalmodels "our-little-chatik/internal/models"
)

type MailerQueue struct {
	ch *amqp.Channel
}

func NewMailerQueue(ch *amqp.Channel) *MailerQueue {
	return &MailerQueue{
		ch: ch,
	}
}

func (m MailerQueue) PutActivationTask(request internalmodels.ActivationTask) internalmodels.StatusCode {
	q, err := m.ch.QueueDeclare(
		"mail", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)

	body, err := json.Marshal(&request)
	if err != nil {
		return internalmodels.InternalError
	}

	err = m.ch.PublishWithContext(
		context.Background(),
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		slog.Error(err.Error())
		return internalmodels.BadRequest
	}
	return internalmodels.OK
}
