package helper

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateRequest(ctx *fiber.Ctx, c *validator.Validate, request any) error {
	err := ctx.BodyParser(request)
	if err != nil {
		return err
	}

	if err := c.Struct(request); err != nil {
		var errMessage string
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
			fieldName := err.Field()
			field, _ := reflect.TypeOf(request).Elem().FieldByName(fieldName)
			jsonField, _ := field.Tag.Lookup("json")
			errMessage = jsonField + " is " + err.ActualTag()
		}

		return fmt.Errorf("%v", errMessage)
	}

	return nil
}
