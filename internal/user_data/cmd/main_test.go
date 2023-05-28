package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/user_data/internal/models"

	"github.com/google/uuid"
)

type testItem struct {
	path         string
	method       string
	prepareBody  func() []byte
	expectedCode int
}

func TestAPI(t *testing.T) {
	tests := []struct {
		name       string
		tasks      []testItem
		testCookie *http.Cookie
		client     http.Client
	}{
		{
			name: "signup login logout logout",
			tasks: []testItem{
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models.UserData{
							User: models2.User{
								Name:     "test",
								Surname:  "test",
								Nickname: "test",
								Password: "test",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusSeeOther,
				},
				{
					path:   "/api/v1/auth/login",
					method: "POST",
					prepareBody: func() []byte {
						person := models.UserData{
							User: models2.User{
								Nickname: "test",
								Password: "test",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusSeeOther,
				},
				{
					path:   "/api/v1/auth/logout",
					method: "DELETE",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusSeeOther,
				},
				{
					path:   "/api/v1/auth/logout",
					method: "DELETE",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusBadRequest,
				},
			},
		},
		{
			name: "signup signup",
			tasks: []testItem{
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models.UserData{
							User: models2.User{
								Name:     "test1",
								Surname:  "test1",
								Nickname: "test1",
								Password: "test1",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusSeeOther,
				},
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models.UserData{
							User: models2.User{
								Name:     "test1",
								Surname:  "test1",
								Nickname: "test2",
								Password: "test1",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusSeeOther,
				},
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models.UserData{
							User: models2.User{
								Name:     "test1",
								Surname:  "test1",
								Nickname: "test1",
								Password: "test1",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusBadRequest,
				},
			},
		},
		{
			name: "signup check person with cookie",
			tasks: []testItem{
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models.UserData{
							User: models2.User{
								Name:     "test3",
								Surname:  "test3",
								Nickname: "test3",
								Password: "test3",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusSeeOther,
				},
				{
					path:   "/api/v1/user/me",
					method: "GET",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusOK,
				},
			},
		},
		{
			name: "signup find another user",
			tasks: []testItem{
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models.UserData{
							User: models2.User{
								Name:     "test4",
								Surname:  "test4",
								Nickname: "test4",
								Password: "test4",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusSeeOther,
				},
				{
					path:   "/api/v1/user/search?name=test1",
					method: "GET",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusOK,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := http.Client{}
			for _, task := range tt.tasks {
				host := os.Getenv("TEST_HOST")
				req, err := http.NewRequest(task.method, host+task.path,
					bytes.NewBuffer(task.prepareBody()))
				if err != nil {
					t.Fatalf("Failed to prepare a request: %s", err)
				}

				if tt.testCookie != nil {
					req.AddCookie(tt.testCookie)
				}

				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
				resp, err := client.Do(req)
				if err != nil {
					t.Fatalf("Failed while performing a request: %s", err)
				}

				if resp.StatusCode != task.expectedCode {
					t.Fatalf("%s: returned status code is wrong %d, expected %d",
						task.path, resp.StatusCode, task.expectedCode)
				}

				if len(resp.Cookies()) != 0 {
					tt.testCookie = resp.Cookies()[0]
				}
			}
		})
	}
}

/* --------------------------------------------------------------------- */
/* --------------------------- BENCH TESTS ----------------------------- */
/* --------------------------------------------------------------------- */

func BenchmarkAPISignUpSeq(b *testing.B) {
	prepareBody := func() []byte {
		person := models.UserData{
			User: models2.User{
				Name:     "test3",
				Surname:  "test3",
				Nickname: uuid.New().String(),
				Password: "test3",
			},
		}
		body, _ := json.Marshal(person)
		return body
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := http.Client{}
		host := os.Getenv("TEST_HOST")
		req, err := http.NewRequest("POST", host+"/api/v1/auth/signup",
			bytes.NewBuffer(prepareBody()))
		if err != nil {
			b.Fatalf("Failed to prepare a request: %s", err)
		}

		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Failed while performing a request: %s", err)
		}

		if resp.StatusCode != http.StatusSeeOther {
			b.Fatalf("returned status code is wrong %d, expected %d",
				resp.StatusCode, http.StatusSeeOther)
		}
	}
}

func BenchmarkAPILoginSeq(b *testing.B) {
	prepareBody := func() []byte {
		person := models.UserData{
			User: models2.User{
				Name:     "test1",
				Surname:  "test1",
				Nickname: "test1",
				Password: "test1",
			},
		}
		body, _ := json.Marshal(person)
		return body
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := http.Client{}
		host := os.Getenv("TEST_HOST")
		req, err := http.NewRequest("POST", host+"/api/v1/auth/login",
			bytes.NewBuffer(prepareBody()))
		if err != nil {
			b.Fatalf("Failed to prepare a request: %s", err)
		}

		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Failed while performing a request: %s", err)
		}

		if resp.StatusCode != http.StatusSeeOther {
			b.Fatalf("returned status code is wrong %d, expected %d",
				resp.StatusCode, http.StatusSeeOther)
		}
	}
}
