package action

type SaveRequest struct {
	RecipeID string `json:"recipe_id"`
}

type LikeRequest struct {
	RecipeID string `json:"recipe_id"`
}
