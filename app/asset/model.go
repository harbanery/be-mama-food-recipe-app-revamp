package asset

import "mime/multipart"

type AssetRequest struct {
	File []*multipart.FileHeader `json:"file"  validate:"required"`
}

type AssetResponse struct {
	FileURL  *string `json:"file_url"`
	FileType *string `json:"file_type"`
}

type DeleteAssetRequest struct {
	FileURL []string `json:"file_url"  validate:"required"`
}
