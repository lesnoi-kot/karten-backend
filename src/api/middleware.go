package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
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
		time.Sleep(1 * time.Second)
		return next(c)
	}
}
