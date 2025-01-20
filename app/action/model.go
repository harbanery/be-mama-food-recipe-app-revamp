package action

type ActionRequest struct {
	RecipeID string `json:"recipe_id" validate:"required"`
}

type CommentRequest struct {
	RecipeID    string `json:"recipe_id" validate:"required"`
	Description string `json:"description" validate:"required"`
}
