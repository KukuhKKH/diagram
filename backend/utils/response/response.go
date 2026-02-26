package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Messages Alias for any slice
type Messages = []any

type Error struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

// error makes it compatible with the error interface
func (e *Error) Error() string {
	return fmt.Sprint(e.Message)
}

// Response A struct to return normal response
type Response struct {
	Code     int      `json:"code"`
	Messages Messages `json:"messages"`
	Data     any      `json:"data,omitempty"`
	Meta     any      `json:"meta,omitempty"`
}

// IsProduction nothiing to describe this fucking variable
var IsProduction bool

// ErrorHandler Default error handler
var ErrorHandler = func(ctx *fiber.Ctx, err error) error {
	resp := Response{
		Code: fiber.StatusInternalServerError,
	}

	// handle errors
	if c, ok := err.(validator.ValidationErrors); ok {
		resp.Code = fiber.StatusUnprocessableEntity
		resp.Messages = Messages{removeTopStruct(c.Translate(trans))}
	} else if c, ok := err.(*fiber.Error); ok {
		resp.Code = c.Code
		resp.Messages = Messages{c.Message}
	} else if c, ok := err.(*Error); ok {
		resp.Code = c.Code
		resp.Messages = Messages{c.Message}

		if resp.Messages == nil {
			resp.Messages = Messages{err}
		}
	} else {
		resp.Messages = Messages{err.Error()}
	}

	if !IsProduction {
		log.Error().Err(err).Msg("From: Fiber's error handler")
	}

	return Resp(ctx, resp)
}

// Resp function to return pretty json response
func Resp(c *fiber.Ctx, resp Response) error {
	if resp.Code == 0 {
		resp.Code = fiber.StatusOK
	}

	c.Status(resp.Code)
	return c.JSON(resp)
}

// remove unecessary fields from validator message
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}

	for field, msg := range fields {
		stripStruct := field[strings.Index(field, ".")+1:]
		res[stripStruct] = msg
	}

	return res
}
