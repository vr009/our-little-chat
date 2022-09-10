package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"our-little-chatik/internal/models"
)

type ChatListRepo struct {
	col *mongo.Collection
}

func NewChatListRepo(col *mongo.Collection) *ChatListRepo {
	return &ChatListRepo{col: col}
}

func (clr *ChatListRepo) FetchChatList(user models.User) ([]models.ChatItem, error) {
	col := clr.col
	ctx := context.TODO()

	findOpts := options.Find().
		SetSort(bson.D{{"last_read", -1}})

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"participant", bson.D{{"$eq", user.UserID}}}},
			}},
	}

	cursor, err := col.Find(ctx, filter, findOpts)
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
