package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"our-little-chatik/internal/models"
)

type AuthHandler struct {
	client http.Client
	cfg    *models.ServiceRouterConfig
}

func NewAuthHandler(client http.Client, cfg *models.ServiceRouterConfig) *AuthHandler {
	return &AuthHandler{client: client, cfg: cfg}
}

func (handler *AuthHandler) AddUser(user models.User) (*models.Session, error) {
	newSession := models.Session{UserID: user.UserID}
	sessionB, err := json.Marshal(newSession)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(sessionB)

	url := handler.cfg.GetPath("AddUser")
	req, err := handler.client.Post(url, "", body)
	if err != nil {
		return nil, err
	}
	if req.StatusCode != 200 {
		return nil, fmt.Errorf("failed to add a user, status code: %d", req.StatusCode)
	}

	session := &models.Session{}
	err = json.NewDecoder(req.Body).Decode(&session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (handler *AuthHandler) RemoveUser(session models.Session) error {
	req, err := http.NewRequest("DELETE", handler.cfg.GetPath("DeleteUser"), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Token", session.Token)
	resp, err := handler.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get a user")
	}

	return nil
}

func (handler *AuthHandler) GetSession(user models.User) (*models.Session, error) {
	req, err := http.NewRequest("GET", handler.cfg.GetPath("GetSession"), nil)
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
	req, err := http.NewRequest("GET", handler.cfg.GetPath("GetUser"), nil)
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