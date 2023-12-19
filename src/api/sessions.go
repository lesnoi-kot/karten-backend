package api

import (
	"errors"
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

const (
	USER_SESSION_KEY    = "session"
	SESSION_KEY_USER_ID = "user-id"
)

func getUserSession(c echo.Context) (*sessions.Session, error) {
	return session.Get(USER_SESSION_KEY, c)
}

func getUserID(c echo.Context) (store.UserID, error) {
	sess, err := getUserSession(c)
	if err != nil {
		return 0, fmt.Errorf("Cannot retrieve session: %w", err)
	}

	userID, ok := sess.Values[SESSION_KEY_USER_ID]
	if !ok {
		return 0, errors.New("Empty session")
	}

	if _, ok = userID.(store.UserID); !ok {
		return 0, errors.New("Invalid session user id")
	}

	return userID.(store.UserID), nil
}

func getUser(c echo.Context) *store.User {
	return c.Get("user").(*store.User)
}

func setUserSession(c echo.Context, userID store.UserID) error {
	sess, _ := session.Get(USER_SESSION_KEY, c)
	sess.Values[SESSION_KEY_USER_ID] = userID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("Session update error: %w", err)
	}

	return nil
}
