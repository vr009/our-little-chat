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
	cfg    *models.ServiceRouterConfig
}

func NewUserDataHandler(client http.Client, cfg *models.ServiceRouterConfig) *UserDataHandler {
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

	resp, err := handler.client.Post(handler.cfg.GetPath("AddUser"), "", body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to add a user")
	}

	newUser := &models.User{}
	err = json.NewDecoder(resp.Body).Decode(&newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (handler *UserDataHandler) RemoveUser(user *models.User) (*models.User, error) {
	req, err := http.NewRequest("DELETE", handler.cfg.GetPath("DeleteUser"), nil)
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

func (handler *UserDataHandler) CheckUser(user *models.User) error {
	userB, err := json.Marshal(user)
	if err != nil {
		return err
	}
	body := bytes.NewReader(userB)

	resp, err := handler.client.Post(handler.cfg.GetPath("CheckUser"), "", body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to add a user")
	}
	return nil
}

func (handler *UserDataHandler) FindUser(name string) ([]models.User, error) {
	req, err := http.NewRequest("GET", handler.cfg.GetPath("FindUser"), nil)
	if err != nil {
		return nil, err
	}
	req.URL.Query().Set("name", name)
	resp, err := handler.client.Do(req)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to find a user")
	}

	var users []models.User
	err = json.NewDecoder(resp.Body).Decode(&users)

	return users, err
}
