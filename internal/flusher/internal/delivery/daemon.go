package delivery

import (
	"context"
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

func (d *FlusherD) Work(ctx context.Context, period int) {
	ticker := time.NewTicker(time.Second * time.Duration(period))
	for {
		select {
		case <-ticker.C:
			messages, err := d.queueRepo.FetchAllMessages()
			if err != nil {
				continue
			}
			err = d.persistantRepo.PersistAllMessages(messages)
			//log.Println("persisted", messages)
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			log.Println("work loop ended")
			return
		}
	}
}
