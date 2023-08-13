package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	models2 "our-little-chatik/internal/models"
)

type testItem struct {
	path          string
	method        string
	prepareBody   func() []byte
	expectedCode  int
	checkResponse func(resp *http.Response) error
}

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	tests := []struct {
		name       string
		tasks      []testItem
		testCookie *http.Cookie
		client     http.Client
	}{
		{
			name: "signup get all users",
			tasks: []testItem{
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models2.UserData{
							User: models2.User{
								Name:     "test5",
								Surname:  "test5",
								Nickname: "test5",
								Password: "test5",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusSeeOther,
					checkResponse: func(resp *http.Response) error {
						return nil
					},
				},
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models2.UserData{
							User: models2.User{
								Name:     "test6",
								Surname:  "test6",
								Nickname: "test6",
								Password: "test6",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusSeeOther,
					checkResponse: func(resp *http.Response) error {
						return nil
					},
				},
				{
					path:   "/api/v1/user/all",
					method: "GET",
					prepareBody: func() []byte {
						return nil
					},

					expectedCode: http.StatusOK,
					checkResponse: func(resp *http.Response) error {
						users := []models2.UserData{}
						err := json.NewDecoder(resp.Body).Decode(&users)
						if len(users) != 2 {
							return fmt.Errorf("not enough users found")
						}
						for _, user := range users {
							if user.Nickname == "" {
								return fmt.Errorf("empty nickname by a found user")
							}
						}
						return err
					},
				},
			},
		},
		{
			name: "signup login logout logout",
			tasks: []testItem{
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models2.UserData{
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
					checkResponse: func(resp *http.Response) error {
						return nil
					},
				},
				{
					path:   "/api/v1/auth/login",
					method: "POST",
					prepareBody: func() []byte {
						person := models2.UserData{
							User: models2.User{
								Nickname: "test",
								Password: "test",
							},
						}
						body, _ := json.Marshal(person)
						return body
					},
					expectedCode: http.StatusSeeOther,
					checkResponse: func(resp *http.Response) error {
						return nil
					},
				},
				{
					path:   "/api/v1/auth/logout",
					method: "DELETE",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusSeeOther,
					checkResponse: func(resp *http.Response) error {
						return nil
					},
				},
				{
					path:   "/api/v1/auth/logout",
					method: "DELETE",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusUnauthorized,
					checkResponse: func(resp *http.Response) error {
						return nil
					},
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
						person := models2.UserData{
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
					checkResponse: func(resp *http.Response) error {
						return nil
					},
				},
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models2.UserData{
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
					checkResponse: func(resp *http.Response) error {
						return nil
					},
				},
				{
					path:   "/api/v1/auth/signup",
					method: "POST",
					prepareBody: func() []byte {
						person := models2.UserData{
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
					checkResponse: func(resp *http.Response) error {
						return nil
					},
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
						person := models2.UserData{
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
					checkResponse: func(resp *http.Response) error {
						return nil
					},
				},
				{
					path:   "/api/v1/user/me",
					method: "GET",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusOK,
					checkResponse: func(resp *http.Response) error {
						return nil
					},
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
						person := models2.UserData{
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
					checkResponse: func(resp *http.Response) error {
						return nil
					},
				},
				{
					path:   "/api/v1/user/search?name=test1",
					method: "GET",
					prepareBody: func() []byte {
						return nil
					},
					expectedCode: http.StatusOK,
					checkResponse: func(resp *http.Response) error {
						person := models2.UserData{}
						err := json.NewDecoder(resp.Body).Decode(&person)
						if err != nil {
							return err
						}
						if person.Nickname != "test4" {
							return fmt.Errorf("wrong person returned: %v", person)
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
				req, err := http.NewRequest(task.method, host+task.path,
					bytes.NewBuffer(task.prepareBody()))
				if err != nil {
					t.Fatalf("Failed to prepare a request: %s", err)
				}
				req.Header.Set("Content-Type", "application/json")

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
					t.Fatalf("%s: returned status code is wrong %d, expected %d test: %s",
						task.path, resp.StatusCode, task.expectedCode, tt.name)
				}

				if len(resp.Cookies()) != 0 {
					tt.testCookie = resp.Cookies()[0]
				}
			}
		})
	}
}

func TestCRUDAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	username := os.Getenv("ADMIN_USER")
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" || password == "" {
		t.Fatalf("ADMIN_USER or ADMIN_PASSWORD are not provided for authentication")
	}

	person := models2.UserData{
		User: models2.User{
			Name:     "test7",
			Surname:  "test7",
			Nickname: "test7",
			Password: "test7",
		},
	}
	body, _ := json.Marshal(person)

	client := http.Client{}
	host := os.Getenv("TEST_HOST")
	req, err := http.NewRequest("POST", host+"/api/v1/auth/signup",
		bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to prepare a request: %s", err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed while performing a request: %s", err)
	}

	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("%s: returned status code is wrong %d, expected %d",
			"/api/v1/auth/signup", resp.StatusCode, http.StatusSeeOther)
	}

	// Create a new user_data
	person = models2.UserData{
		User: models2.User{
			Name:     "test8",
			Surname:  "test8",
			Nickname: "test8",
			Password: "test8",
		},
	}
	body, _ = json.Marshal(person)

	req, err = http.NewRequest("POST", host+"/api/v1/admin/user",
		bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to prepare a request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed while performing a request: %s", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("wrong status code expected %d actual %d",
			http.StatusCreated, resp.StatusCode)
	}

	var createdPerson models2.UserData
	err = json.NewDecoder(resp.Body).Decode(&createdPerson)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Get
	req, err = http.NewRequest("GET", host+"/api/v1/admin/user?id="+createdPerson.UserID.String(),
		nil)
	if err != nil {
		t.Fatalf("Failed to prepare a request: %s", err)
	}
	req.SetBasicAuth(username, password)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed while performing a request: %s", err)
	}

	var foundPerson models2.UserData
	err = json.NewDecoder(resp.Body).Decode(&foundPerson)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if foundPerson.UserID != createdPerson.UserID {
		t.Fatalf("the id of found user_data is not the same as returned"+
			" after its creation: %s %s", foundPerson.UserID.String(), createdPerson.UserID.String())
	}
}

/* --------------------------------------------------------------------- */
/* --------------------------- BENCH TESTS ----------------------------- */
/* --------------------------------------------------------------------- */

func BenchmarkAPISignUpSeq(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}
	prepareBody := func() []byte {
		person := models2.UserData{
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
		req.Header.Set("Content-Type", "application/json")
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
	if testing.Short() {
		b.Skip()
	}
	prepareBody := func() []byte {
		person := models2.UserData{
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
		req.Header.Set("Content-Type", "application/json")

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

func BenchmarkTestAPIParallel(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}
	prepareBody := func() []byte {
		person := models2.UserData{
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

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client := http.Client{}
			host := os.Getenv("TEST_HOST")
			req, err := http.NewRequest("POST", host+"/api/v1/auth/login",
				bytes.NewBuffer(prepareBody()))
			if err != nil {
				b.Fatalf("Failed to prepare a request: %s", err)
			}
			req.Header.Set("Content-Type", "application/json")

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
	})
}
