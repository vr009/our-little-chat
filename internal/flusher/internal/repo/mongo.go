package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"our-little-chatik/internal/models"
	"time"
)

type MongoRepo struct {
	conn *mongo.Client
}

func NewMongoRepo(conn *mongo.Client) *MongoRepo {
	return &MongoRepo{conn: conn}
}

func (repo *MongoRepo) PersistAll(msgs []models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	docs := []interface{}{}
	for _, msg := range msgs {
		docs = append(docs, msg)
	}

	collection := repo.conn.Database("messages").Collection("message_list")
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	fmt.Println("persisted:", msgs)
	return nil
}
