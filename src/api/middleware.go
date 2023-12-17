package api

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

var requireId echo.MiddlewareFunc = requireParam("id")

func requireIntParam(name string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			param, err := strconv.Atoi(c.Param(name))
			if err != nil {
				return echo.ErrBadRequest
			}

			c.Set(name, param)
			return next(c)
		}
	}
}

func requireParam(name string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Param(name) == "" {
				return echo.ErrBadRequest
			}

			return next(c)
		}
	}
}

type EchoValidator struct {
	validator *validator.Validate
}

func newEchoValidator() EchoValidator {
	return EchoValidator{validator.New()}
}

func (v EchoValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func parseError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}

		return err
	}
}

func emulateDelay(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		time.Sleep(200 * time.Millisecond)
		return next(c)
	}
}

func (service *APIService) makeRequireAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := getUserID(c)
			if err != nil {
				return echo.ErrUnauthorized
			}
			c.Set("userID", userID)
			return next(c)
		}
	}
}

func (service *APIService) makeInjectUserMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userService, err := service.getUserService(c)
			if err != nil {
				return echo.ErrUnauthorized
			}

			user, err := userService.GetUser(&entityservices.GetUserOptions{
				FullInfo:      true,
				IncludeAvatar: true,
			})

			if errors.Is(err, store.ErrNotFound) {
				return echo.ErrUnauthorized
			}
			if err != nil {
				return err
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
