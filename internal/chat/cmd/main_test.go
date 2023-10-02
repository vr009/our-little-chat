package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"sort"
	"testing"
	"time"
)

type testItem struct {
	path          string
	method        string
	prepareBody   func() []byte
	expectedCode  int
	preparePath   func(str string) string
	checkResponse func(resp *http.Response) error
}

func TestChatAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var testChat models2.Chat
	token, _ := pkg.GenerateJWTTokenV2(models2.User{
		ID: uuid.New(),
	}, false)

	testCookie := &http.Cookie{Name: "Token", Value: token, Path: "/"}

	testChat = models2.Chat{
		Name: "test",
		Participants: []uuid.UUID{
			uuid.New(),
		},
		PhotoURL: "test",
	}

	tests := []struct {
		name       string
		tasks      []testItem
		testCookie *http.Cookie
		client     http.Client
	}{
		{
			testCookie: testCookie,
			name:       "create chat get chat info",
			tasks: []testItem{
				{
					path: "/api/v1/new",
					preparePath: func(str string) string {
						return "/api/v1/new"
					},
					method: "POST",
					prepareBody: func() []byte {
						chat := testChat
						body, _ := json.Marshal(chat)
						return body
					},
					expectedCode: http.StatusCreated,
					checkResponse: func(resp *http.Response) error {
						chat := models2.Chat{}
						err := json.NewDecoder(resp.Body).Decode(&chat)
						if err != nil {
							return err
						}
						t.Log(chat.ChatID.String())
						testChat.ChatID = chat.ChatID
						testChat.Participants = chat.Participants
						return nil
					},
				},
				{
					path: "/api/v1/chat?chat_id=",
					preparePath: func(str string) string {
						return "/api/v1/chat?chat_id=" + str
					},
					method: "GET",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusOK,
					checkResponse: func(resp *http.Response) error {
						chat := models2.Chat{}
						err := json.NewDecoder(resp.Body).Decode(&chat)
						if err != nil {
							return err
						}
						if chat.ChatID != testChat.ChatID {
							return fmt.Errorf("wrong chat id")
						}
						if chat.PhotoURL != testChat.PhotoURL {
							return fmt.Errorf("wrong photo url")
						}
						if len(chat.Participants) != 2 {
							return fmt.Errorf("wrong photo url")
						}
						return nil
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := http.Client{}
			for _, task := range tt.tasks {
				host := os.Getenv("TEST_HOST")
				req, err := http.NewRequest(task.method, host+task.preparePath(testChat.ChatID.String()),
					bytes.NewBuffer(task.prepareBody()))
				if err != nil {
					t.Fatalf("Failed to prepare a request: %s", err)
				}

				req.AddCookie(tt.testCookie)
				req.Header.Set("Content-Type", "application/json")

				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
				resp, err := client.Do(req)
				if err != nil {
					t.Fatalf("Failed while performing a request: %s", err)
				}

				if resp.StatusCode != task.expectedCode {
					t.Fatalf("%s: returned status code is wrong %d, expected %d test: %s",
						task.path, resp.StatusCode, task.expectedCode, tt.name)
				}
				if err := task.checkResponse(resp); err != nil {
					t.Fatalf(err.Error())
				}
			}
		})
	}
}

func putMessagesToRedis(msgs models2.Messages, cl *redis.Client) error {
	for _, msg := range msgs {
		bData, _ := json.Marshal(&msg)
		cl.Set(context.Background(),
			fmt.Sprintf("%s_%s", msg.ChatID, msg.MsgID), bData, time.Minute*5)
	}
	return nil
}

func putMessagesToPostgres(msgs models2.Messages, conn *pgxpool.Pool) error {
	queryStr := "INSERT INTO messages VALUES ($1, $2, $3, $4, $5)"
	for _, msg := range msgs {
		_, err := conn.Exec(context.Background(), queryStr, msg.MsgID, msg.ChatID, msg.SenderID, msg.Payload, msg.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestMsgAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testChat := models2.Chat{ChatID: uuid.New()}
	testSenderID := uuid.New()
	userID := uuid.New()

	msg1 := models2.Message{
		MsgID:     uuid.New(),
		ChatID:    testChat.ChatID,
		Payload:   "test1",
		SenderID:  testSenderID,
		CreatedAt: 100,
	}
	msg2 := models2.Message{
		MsgID:     uuid.New(),
		ChatID:    testChat.ChatID,
		Payload:   "test2",
		SenderID:  userID,
		CreatedAt: 1000,
	}

	msg3 := models2.Message{
		MsgID:     uuid.New(),
		ChatID:    testChat.ChatID,
		Payload:   "test3",
		SenderID:  testSenderID,
		CreatedAt: 12,
	}
	msg4 := models2.Message{
		MsgID:     uuid.New(),
		ChatID:    testChat.ChatID,
		Payload:   "test4",
		SenderID:  userID,
		CreatedAt: 13,
	}

	connStr, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		t.Fatalf("connection string not found")
	}
	testPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
	defer testPool.Close()

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		t.Fatalf("empty redis host")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		t.Fatalf("empty redis port")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		t.Fatalf("empty redis password")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
	})

	err = putMessagesToRedis(models2.Messages{msg1, msg2}, redisClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = putMessagesToPostgres(models2.Messages{msg3, msg4}, testPool)
	if err != nil {
		t.Fatalf(err.Error())
	}

	token, _ := pkg.GenerateJWTTokenV2(models2.User{
		ID: userID,
	}, false)

	testCookie := &http.Cookie{Name: "Token", Value: token, Path: "/"}

	tests := []struct {
		name       string
		tasks      []testItem
		testCookie *http.Cookie
		client     http.Client
	}{
		{
			testCookie: testCookie,
			name:       "conv",
			tasks: []testItem{
				{
					path: "/api/v1/conv",
					preparePath: func(chatID string) string {
						return "/api/v1/conv?limit=10&offset=0&chat_id=" + chatID
					},
					method: "GET",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusOK,
					checkResponse: func(resp *http.Response) error {
						msgs := models2.Messages{}
						err := json.NewDecoder(resp.Body).Decode(&msgs)
						if err != nil {
							return err
						}
						if len(msgs) != 4 {
							return fmt.Errorf("wrong len of msg slice")
						}
						if !sort.IsSorted(msgs) {
							return fmt.Errorf("msgs are not sorted")
						}
						return nil
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := http.Client{}
			for _, task := range tt.tasks {
				host := os.Getenv("TEST_HOST")
				req, err := http.NewRequest(task.method, host+task.preparePath(testChat.ChatID.String()),
					bytes.NewBuffer(task.prepareBody()))
				if err != nil {
					t.Fatalf("Failed to prepare a request: %s", err)
				}

				req.AddCookie(tt.testCookie)

				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
				resp, err := client.Do(req)
				if err != nil {
					t.Fatalf("Failed while performing a request: %s", err)
				}

				if resp.StatusCode != task.expectedCode {
					t.Fatalf("%s: returned status code is wrong %d, expected %d test: %s",
						task.path, resp.StatusCode, task.expectedCode, tt.name)
				}
				if err := task.checkResponse(resp); err != nil {
					t.Fatalf(err.Error())
				}
			}
		})
	}
}
