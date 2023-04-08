package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/api/oauth"
	"github.com/lesnoi-kot/karten-backend/src/settings"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

var githubOAuthProvider oauth.GitHubProvider

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

	db_social_id := fmt.Sprintf("%s_%s", oauthProvider.GetName(), userInfo.ID)
	db_user, err := api.store.Users.GetBySocialID(context.Background(), db_social_id)
	if errors.Is(err, store.ErrNotFound) {
		db_user = &store.User{
			SocialID: db_social_id,
			Name:     userInfo.Name,
			Login:    userInfo.Login,
			Email:    userInfo.Email,
			URL:      userInfo.URL,
		}

		if err := api.store.Users.Add(context.Background(), db_user); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	sess, _ := session.Get(USER_SESSION_KEY, c)
	sess.Values[SESSION_KEY_USER_ID] = db_user.ID
	sess.Options = &sessions.Options{
		Path:     "/",
		Domain:   settings.AppConfig.CookieDomain,
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("Session update rror: %w", err)
	}

	api.logger.Debugf("session: %v, user: %v", sess.Values, userInfo)
	return c.Redirect(http.StatusTemporaryRedirect, settings.AppConfig.FrontendURL)
}
