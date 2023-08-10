package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"testing"
)

type testItem struct {
	path          string
	method        string
	prepareBody   func() []byte
	expectedCode  int
	preparePath   func(str string) string
	checkResponse func(resp *http.Response) error
}

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var testChat models2.Chat
	token, _ := pkg.GenerateJWTToken(models2.User{
		UserID: uuid.New(),
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
