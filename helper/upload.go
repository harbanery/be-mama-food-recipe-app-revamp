package helper

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
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

func DeleteFile(context *context.Context, cloudinary *cloudinary.Cloudinary, urlString *string) error {
	filenameWithExtension := path.Base(*urlString)

	decodedFilenameWithExtension, err := url.QueryUnescape(filenameWithExtension)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to decode url: "+err.Error())
	}

	filename := strings.TrimSuffix(decodedFilenameWithExtension, path.Ext(decodedFilenameWithExtension))

	destroyParams := uploader.DestroyParams{
		PublicID: filename,
	}

	if _, err = cloudinary.Upload.Destroy(*context, destroyParams); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
