package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lesnoi-kot/karten-backend/src/settings"
)

type GitHubProvider struct{}

func (p GitHubProvider) GetName() string {
	return "github"
}

func (p GitHubProvider) GetAccessToken(c *http.Client, code string) (string, error) {
	params := url.Values{
		"client_id":     {settings.AppConfig.GithubClientID},
		"client_secret": {settings.AppConfig.GithubClientSecret},
		"code":          {code},
	}

	resp, err := c.PostForm("https://github.com/login/oauth/access_token", params)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		return "", errors.New("GitHub OAuth /login/oauth/access_token returned non 2xx code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	auth_resp, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}

	auth_error := auth_resp.Get("error")

	if auth_error != "" {
		return "", errors.New(auth_error)
	}

	return auth_resp.Get("access_token"), nil
}

func (p GitHubProvider) GetUser(c *http.Client, accessToken string) (*UserInfo, error) {
	user_request, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	user_request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	user_resp, err := c.Do(user_request)
	if err != nil {
		return nil, err
	}

	defer user_resp.Body.Close()
	user_body, err := ioutil.ReadAll(user_resp.Body)
	if err != nil {
		return nil, err
	}

	type GitHubUser struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		URL       string `json:"html_url"`
		AvatarURL string `json:"avatar_url"`
	}

	user := new(GitHubUser)
	if err := json.Unmarshal(user_body, user); err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:    fmt.Sprint(user.ID),
		Name:  user.Name,
		Login: user.Login,
		Email: user.Email,
		URL:   user.URL,
	}, nil
}
