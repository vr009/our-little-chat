package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	"our-little-chatik/internal/models"

	"github.com/golang/glog"
)

func AuthHook(r *http.Request, authUrl string) (*models.User, error) {
	session := models.Session{}

	cookie, err := r.Cookie("Token")
	session.Token = cookie.Value
	if session.Token == "" {
		err = fmt.Errorf("no cookie provided")
		return nil, err
	}

	cl := http.Client{}
	req, err := http.NewRequest(http.MethodGet, authUrl, nil)
	if err != nil {
		err := fmt.Errorf("failed to build a request %s: %s", authUrl, err.Error())
		glog.Error(err.Error())
		return nil, err
	}

	req.Header.Set("Token", session.Token)
	resp, err := cl.Do(req)
	if err != nil {
		err := fmt.Errorf("failed to hook %s: %s", authUrl, err.Error())
		glog.Error(err.Error())
		return nil, err
	}
	glog.Warningf("here")
	user := models.User{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		err := fmt.Errorf("failed to decode an answer from %s: %s", authUrl, err.Error())
		glog.Error(err.Error())
		return nil, err
	}
	glog.Warningf("here %v", user)

	return &user, nil
}
