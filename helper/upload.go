package helper

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
)

func UploadFile(context *context.Context, cloudinary *cloudinary.Cloudinary, file *multipart.FileHeader) (*string, error) {
	var url *string

	src, err := file.Open()
	if err != nil {
		return url, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	fileNameWithoutExt := file.Filename[:len(file.Filename)-len(ext)]

	uploadParams := uploader.UploadParams{
		PublicID:  fmt.Sprintf("%d_%s", time.Now().Unix(), fileNameWithoutExt),
		Overwrite: api.Bool(true),
	}

	uploadResult, err := cloudinary.Upload.Upload(*context, src, uploadParams)
	if err != nil {
		return url, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &uploadResult.SecureURL, nil
}

func DeleteFile(context *context.Context, cloudinary *cloudinary.Cloudinary, url *string) error {
	filenameWithExtension := path.Base(*url)

	filename := strings.TrimSuffix(filenameWithExtension, path.Ext(filenameWithExtension))

	destroyParams := uploader.DestroyParams{
		PublicID: filename,
	}

	_, err := cloudinary.Upload.Destroy(*context, destroyParams)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
