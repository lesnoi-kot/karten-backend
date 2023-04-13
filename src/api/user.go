package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/authservice"
	"github.com/lesnoi-kot/karten-backend/src/authservice/oauth"
	"github.com/lesnoi-kot/karten-backend/src/settings"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

var githubOAuthProvider oauth.GitHubProvider

type UserDTO struct {
	ID          int       `json:"id"`
	SocialID    string    `json:"social_id"`
	Name        string    `json:"name"`
	Login       string    `json:"login"`
	Email       string    `json:"email"`
	URL         string    `json:"url"`
	DateCreated time.Time `json:"date_created"`
}

func (api *APIService) oauthCallback(c echo.Context) error {
	oauth_code := c.QueryParam("code")
	oauth_authorizer := c.QueryParam("authorizer")

	if oauth_code == "" {
		return c.Redirect(http.StatusTemporaryRedirect, settings.AppConfig.FrontendURL)
	}

	var oauthProvider oauth.OAuthProvider = nil

	switch oauth_authorizer {
	case "github":
		oauthProvider = githubOAuthProvider
	default:
		oauthProvider = githubOAuthProvider
	}

	accessToken, err := oauthProvider.GetAccessToken(http.DefaultClient, oauth_code)
	if err != nil {
		return fmt.Errorf("GetAccessToken error: %w", err)
	}

	userInfo, err := oauthProvider.GetUser(http.DefaultClient, accessToken)
	if err != nil {
		return fmt.Errorf("GetUser error: %w", err)
	}

	authService := authservice.AuthService{Store: api.store}
	db_user, err := authService.Authenticate(userInfo)
	if err != nil {
		return err
	}

	if err := setUserSession(c, db_user.ID); err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, settings.AppConfig.FrontendURL)
}

func (api *APIService) getCurrentUser(c echo.Context) error {
	user, ok := c.Get("user").(*store.User)
	if !ok {
		return echo.ErrUnauthorized
	}
	return c.JSON(http.StatusOK, OK(userToDTO(user)))
}

func (api *APIService) logOut(c echo.Context) error {
	sess, err := getUserSession(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, Error("Cannot retrieve session"))
	}

	sess.Options.MaxAge = -1
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("Session update error: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

func (api *APIService) guestLogIn(c echo.Context) error {
	if err := setUserSession(c, store.GuestUserID); err != nil {
		return err
	}

	user, err := api.store.Users.Get(context.Background(), store.GuestUserID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(userToDTO(user)))
}
