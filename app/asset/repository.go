package asset

type AssetRepository interface{}

type assetRepository struct{}

func NewAssetRepository() AssetRepository {
	return &assetRepository{}
}
