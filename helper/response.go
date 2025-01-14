package helper

import (
	"github.com/gofiber/fiber/v2"
)

type ResponseMessageLog struct {
	Message string `json:"message"`
}

func Response(ctx *fiber.Ctx, code int, message string, data any) *WebResponse[interface{}] {
	res := new(WebResponse[interface{}])
	msg := message

	if ctx.OriginalURL() == "" {
		code = 500
		msg = "Internal Server Error"
		data = EmptyObject()
	}

	res.Code = code
	res.Message = msg
	if data == nil {
		res.Data = nil
	} else {
		res.Data = &data
	}

	return res
}
