package oauth

import "net/http"

type UserInfo struct {
	ID    string
	Name  string
	Login string
	Email string
	URL   string
}

type OAuthProvider interface {
	GetName() string
	GetAccessToken(c *http.Client, code string) (string, error)
	GetUser(c *http.Client, accessToken string) (*UserInfo, error)
}
