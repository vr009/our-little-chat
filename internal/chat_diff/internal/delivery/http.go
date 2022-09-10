package delivery

import "net/url"

type AuthDialer struct {
	url url.URL
}

func NewAuthDialer(addr string) *AuthDialer {
	u := url.URL{Path: addr}
	return &AuthDialer{url: u}
}
