package api

import (
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
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

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
