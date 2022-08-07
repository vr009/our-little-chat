package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"our-little-chatik/internal/models"
)

type MongoRepo struct {
	db *mongo.Database
}

func NewMongoRepo(db *mongo.Database) *MongoRepo {
	return &MongoRepo{
		db: db,
	}
}

func (repo *MongoRepo) InsertMessages(msgs []models.Message) error {
	batch := []interface{}{}
	for _, msg := range msgs {
		batch = append(batch, msg)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := repo.db.Collection("messages").InsertMany(ctx, batch)
	if err != nil {
		return err
	}
	if len(res.InsertedIDs) < 1 {
		return fmt.Errorf("nothing inserted")
	}
	return nil
}

func (repo *MongoRepo) GetChat(chat models.Chat, opts models.Opts) ([]models.Message, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	coll := repo.db.Collection("messages")

	findOpts := options.Find().
		SetLimit(opts.Limit).
		SetSkip(opts.Page).
		SetSort(bson.D{{"created_at", -1}})

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"chat_id", bson.D{{"$eq", chat.ChatID}}}},
			}},
	}

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}

	var results []models.Message
	for cursor.Next(ctx) {
		msg := models.Message{}
		err = cursor.Decode(&msg)
		if err != nil {
			return nil, err
		}
		results = append(results, msg)
	}
	if err != nil {
		return nil, err
	}

	return results, nil
}
