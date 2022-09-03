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
			chats, err := d.queueRepo.FetchAllChats()
			if err != nil {
				log.Println(err)
			}
			err = d.persistantRepo.PersistAllChats(chats)
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			log.Println("work loop ended")
			return
		}
	}
}
