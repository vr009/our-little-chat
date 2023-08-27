package delivery

import (
	"encoding/json"
	"net/http"
	"os"
	"our-little-chatik/internal/models"
)

type UserDataClient struct {
	cl      http.Client
	baseURl string
}

const getUserPath = "/api/v1/admin/user"

func NewUserDataClient(cl http.Client, baseURl string) *UserDataClient {
	return &UserDataClient{
		cl:      cl,
		baseURl: baseURl,
	}
}

func (cl UserDataClient) GetUser(user models.UserData) (models.UserData, error) {
	req, err := http.NewRequest("GET",
		cl.baseURl+getUserPath+"?id="+user.UserID.String(), nil)
	req.SetBasicAuth(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"))

	resp, err := cl.cl.Do(req)
	if err != nil {
		return models.UserData{}, err
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return models.UserData{}, err
	}
	return user, nil
}
