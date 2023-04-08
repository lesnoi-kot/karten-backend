package api

import (
	"context"
	"errors"
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
	sess, err := session.Get(USER_SESSION_KEY, c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, Error("Cannot retrieve session"))
	}

	userID, ok := sess.Values[SESSION_KEY_USER_ID]
	if !ok {
		return c.JSON(http.StatusUnauthorized, Error("Empty session"))
	}

	if _, ok = userID.(int); !ok {
		return c.JSON(http.StatusBadRequest, Error("Invalid session"))
	}

	user, err := api.store.Users.Get(context.Background(), userID.(int))

	if errors.Is(err, store.ErrNotFound) {
		return c.JSON(http.StatusUnauthorized, Error("User not found"))
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(userToDTO(user)))
}
