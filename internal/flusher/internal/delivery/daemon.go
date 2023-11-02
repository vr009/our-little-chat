package delivery

import (
	"context"
	"golang.org/x/exp/slog"
	"log"
	"time"

	"our-little-chatik/internal/flusher/internal"
)

type FlusherD struct {
	queueRepo      internal.QueueRepo
	persistantRepo internal.PersistantRepo
}

func NewFlusherD(queueRepo internal.QueueRepo, persistantRepo internal.PersistantRepo) *FlusherD {
	return &FlusherD{queueRepo: queueRepo, persistantRepo: persistantRepo}
}

func (d *FlusherD) Work(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(time.Minute * period)
	for {
		select {
		case <-ticker.C:
			messages, err := d.queueRepo.FetchAllMessages()
			if err != nil {
				continue
			}
			err = d.persistantRepo.PersistAllMessages(messages)
			slog.Info("persisted messages", messages)
			if err != nil {
				log.Println(err)
			}

			lastMessages, err := d.queueRepo.FetchAllLastMessagesOfChats()
			if err != nil {
				log.Println(err)
				continue
			}
			err = d.persistantRepo.PersistAllLastChatMessages(lastMessages)
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			log.Println("work loop ended")
			return
		}
	}
}
