package helper

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"slices"
	"strings"

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

func ValidateImageRequest(file *multipart.FileHeader) error {
	fileSizeLimit := 2000 * 1024
	if file.Size > int64(fileSizeLimit) {
		return fiber.NewError(418, "file size exceeds 2MB limit")
	}

	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png"}
	if !slices.Contains(allowedTypes, file.Header.Get("Content-Type")) {
		return fiber.NewError(419, "invalid type file")
	}

	allowedExtension := []string{"jpeg", "jpg", "png"}
	splitFileName := strings.Split(file.Filename, ".")
	extension := splitFileName[len(splitFileName)-1]
	if !slices.Contains(allowedExtension, extension) {
		return fiber.NewError(419, "invalid type file")
	}

	return nil

}
