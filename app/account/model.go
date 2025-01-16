package account

import "mime/multipart"

type ProfileRequest struct {
	Fullname string `json:"fullname" validate:"required"`
}

type PhotoRequest struct {
	Photo []*multipart.FileHeader `json:"photo" validate:"required"`
}

type ProfileResponse struct {
	ID           string        `json:"id"`
	Fullname     string        `json:"fullname"`
	Photo        string        `json:"photo"`
	MyRecipes    []*Recipe     `json:"my_recipes" gorm:"foreignKey:AuthorID"`
	SavedRecipes []*SaveRecipe `json:"saved_recipes" gorm:"foreignKey:UserID"`
	LikedRecipes []*LikeRecipe `json:"liked_recipes" gorm:"foreignKey:UserID"`
}

type Recipe struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	Image    string `json:"image"`
	AuthorID string `json:"author_id"`
}

type SaveRecipe struct {
	UserID   string  `json:"user_id"`
	RecipeID string  `json:"recipe_id"`
	Recipe   *Recipe `json:"recipe" gorm:"foreignKey:RecipeID"`
}

type LikeRecipe struct {
	UserID   string  `json:"user_id"`
	RecipeID string  `json:"recipe_id"`
	Recipe   *Recipe `json:"recipe" gorm:"foreignKey:RecipeID"`
}
