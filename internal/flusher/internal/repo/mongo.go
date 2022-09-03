package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"our-little-chatik/internal/models"
	"time"
)

type MongoRepo struct {
	msgsCol     *mongo.Collection
	chatListCol *mongo.Collection
}

func NewMongoRepo(msgsCol *mongo.Collection, chatListCol *mongo.Collection) *MongoRepo {
	return &MongoRepo{msgsCol: msgsCol, chatListCol: chatListCol}
}

func (repo *MongoRepo) PersistAllMessages(msgs []models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	docs := []interface{}{}
	for _, msg := range msgs {
		docs = append(docs, msg)
	}

	_, err := repo.msgsCol.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	//fmt.Println("persisted:", msgs)
	return nil
}

func (repo *MongoRepo) PersistAllChats(chats []models.Chat) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	docs := []interface{}{}
	for _, chat := range chats {
		docs = append(docs, chat)
	}

	_, err := repo.chatListCol.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	fmt.Println("persisted chats:", chats)
	return nil
}
