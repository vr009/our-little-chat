package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"our-little-chatik/internal/models"
)

type UserDataHandler struct {
	client http.Client
	cfg    UserDataConfig
}

func NewUserDataHandler(client http.Client, cfg UserDataConfig) *UserDataHandler {
	return &UserDataHandler{client: client, cfg: cfg}
}

type UserDataConfig struct {
	BaseUrl string
	Router  map[string]string
}

func (cfg UserDataConfig) GetPath(method string) string {
	return filepath.Join(cfg.BaseUrl, cfg.Router[method])
}

func (handler *UserDataHandler) AddUser(user *models.User) (*models.User, error) {
	userB, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(userB)

	req, err := handler.client.Post(handler.cfg.GetPath("addUser"), "", body)
	if err != nil {
		return nil, err
	}
	if req.StatusCode != 200 {
		return nil, fmt.Errorf("failed to add a user")
	}

	newUser := &models.User{}
	err = json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (handler *UserDataHandler) RemoveUser(user *models.User) (*models.User, error) {
	req, err := http.NewRequest("DELETE", handler.cfg.GetPath("removeUser"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("UserID", user.UserID.String())
	resp, err := handler.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to remove a user")
	}

	err = json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (handler *UserDataHandler) GetUser(user *models.User) (*models.User, error) {
	req, err := http.NewRequest("GET", handler.cfg.GetPath("getUser"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("UserID", user.UserID.String())
	resp, err := handler.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get a user")
	}

	err = json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
