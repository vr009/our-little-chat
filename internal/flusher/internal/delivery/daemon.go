package delivery

import (
	"context"
	"log"
	"our-little-chatik/internal/flusher/internal"
	"time"
)

type FlusherD struct {
	queueRepo      internal.QueueRepo
	persistantRepo internal.PersistantRepo
}

func NewFlusherD(queueRepo internal.QueueRepo, persistantRepo internal.PersistantRepo) *FlusherD {
	return &FlusherD{queueRepo: queueRepo, persistantRepo: persistantRepo}
}

func (d *FlusherD) Work(ctx context.Context, period int) {
	ticker := time.NewTicker(time.Hour * time.Duration(period))
	for {
		select {
		case <-ticker.C:
			messages, err := d.queueRepo.FetchAll()
			if err != nil {
				log.Println(err)
			}
			err = d.persistantRepo.PersistAll(messages)
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			log.Println("work loop ended")
			return
		}
	}
}
