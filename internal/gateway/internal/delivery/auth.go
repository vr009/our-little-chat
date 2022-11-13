package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"our-little-chatik/internal/models"
)

type AuthHandler struct {
	client http.Client
	cfg    AuthServiceConfig
}

type AuthServiceConfig struct {
	BaseUrl string
	Router  map[string]string
}

func (cfg AuthServiceConfig) GetPath(method string) string {
	return filepath.Join(cfg.BaseUrl, cfg.Router[method])
}

func NewAuthHandler(client http.Client, cfg AuthServiceConfig) *AuthHandler {
	return &AuthHandler{client: client, cfg: cfg}
}

func (handler *AuthHandler) AddUser(user models.User) (*models.Session, error) {
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

	session := &models.Session{}
	err = json.NewDecoder(req.Body).Decode(&session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (handler *AuthHandler) RemoveUser(user models.User) (*models.Session, error) {
	req, err := http.NewRequest("DELETE", handler.cfg.GetPath("deleteUser"), nil)
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

	session := &models.Session{}
	err = json.NewDecoder(req.Body).Decode(&session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (handler *AuthHandler) GetSession(user models.User) (*models.Session, error) {
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

	session := &models.Session{}
	err = json.NewDecoder(req.Body).Decode(&session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (handler *AuthHandler) GetUser(session models.Session) (*models.User, error) {
	req, err := http.NewRequest("GET", handler.cfg.GetPath("getUser"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Token", session.Token)
	resp, err := handler.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get a user")
	}

	user := &models.User{}
	err = json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
