package repo

import (
	"context"
	"time"

	models2 "our-little-chatik/internal/chat/models"
	"our-little-chatik/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	db         *mongo.Database
	chatListDB *mongo.Database
}

func NewMongoRepo(db *mongo.Database, chatListDB *mongo.Database) *MongoRepo {
	return &MongoRepo{
		db:         db,
		chatListDB: chatListDB,
	}
}

func (repo *MongoRepo) GetChatMessages(chat models2.Chat, opts models.Opts) ([]models.Message, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	coll := repo.db.Collection("chat")

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

func (clr *MongoRepo) FetchChatList(user models.User) ([]models.ChatItem, error) {
	col := clr.chatListDB.Collection("chat_list")
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))

	findOpts := options.Find().
		SetSort(bson.D{{"last_read", -1}})

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"participant", bson.D{{"$eq", user.UserID}}}},
			}},
	}

	cursor, err := col.Find(ctx, filter, findOpts)
	if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
		return []models.ChatItem{}, nil
	}
	if err != nil {
		return nil, err
	}

	chats := []models.ChatItem{}
	err = cursor.All(ctx, &chats)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (clr *MongoRepo) InsertChat(chat models2.Chat) error {
	_, err := clr.chatListDB.Collection("chat_list").InsertOne(context.Background(), &chat)
	return err
}
