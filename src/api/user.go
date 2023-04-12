package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type UserDTO struct {
	ID          int       `json:"id"`
	SocialID    string    `json:"social_id"`
	Name        string    `json:"name"`
	Login       string    `json:"login"`
	Email       string    `json:"email"`
	URL         string    `json:"url"`
	DateCreated time.Time `json:"date_created"`
}

func (api *APIService) getCurrentUser(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, Error(err.Error()))
	}

	user, err := api.store.Users.Get(context.Background(), userID)

	if errors.Is(err, store.ErrNotFound) {
		return c.JSON(http.StatusUnauthorized, Error("User not found"))
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(userToDTO(user)))
}

func (api *APIService) logOut(c echo.Context) error {
	sess, err := session.Get(USER_SESSION_KEY, c)
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
