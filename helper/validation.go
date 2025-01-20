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

func ValidateRequest(ctx *fiber.Ctx, valid *validator.Validate, request any) error {
	err := ctx.BodyParser(request)
	if err != nil {
		return err
	}

	if err := valid.Struct(request); err != nil {
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

func ValidateFormRequest(ctx *fiber.Ctx, valid *validator.Validate, request any) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return fmt.Errorf("invalid form-data: %v", err)
	}

	requestValue := reflect.ValueOf(request).Elem()
	requestType := requestValue.Type()

	for i := 0; i < requestType.NumField(); i++ {
		field := requestType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}

		if field.Type == reflect.TypeOf([]*multipart.FileHeader{}) {
			files := form.File[jsonTag]
			if len(files) > 0 {
				requestValue.Field(i).Set(reflect.ValueOf(files))
			}
		} else if field.Type == reflect.TypeOf([]string{}) {
			formValue := form.Value[jsonTag]
			if len(formValue) > 0 {
				requestValue.Field(i).Set(reflect.ValueOf(formValue))
			}
		} else {
			formValue := form.Value[jsonTag]
			if len(formValue) > 0 {
				requestValue.Field(i).SetString(formValue[0])
			}
		}
	}

	if err := valid.Struct(request); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return fmt.Errorf("%v", err)
		}

		var errorMessages []string
		for _, err := range validationErrors {
			fieldName := err.Field()
			if field, ok := requestType.FieldByName(fieldName); ok {
				if jsonField, ok := field.Tag.Lookup("json"); ok {
					fieldName = jsonField
				}
			}
			errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", fieldName, err.ActualTag()))
		}

		return fmt.Errorf("%v", strings.Join(errorMessages, ", "))
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

func ValidateVideoRequest(file *multipart.FileHeader) error {
	allowedTypes := []string{"video/mp4", "video/webm", "video/ogg"}
	if !slices.Contains(allowedTypes, file.Header.Get("Content-Type")) {
		return fiber.NewError(419, "invalid type file")
	}

	allowedExtension := []string{"mp4", "webm", "ogg"}
	splitFileName := strings.Split(file.Filename, ".")
	extension := splitFileName[len(splitFileName)-1]
	if !slices.Contains(allowedExtension, extension) {
		return fiber.NewError(419, "invalid type file")
	}

	return nil
}
