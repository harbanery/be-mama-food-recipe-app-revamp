package account

type ProfileRequest struct {
	Fullname string `json:"fullname" validate:"required"`
}

type PhotoRequest struct {
	Photo string `json:"photo" validate:"required"`
}

type Profile struct {
	ID           string           `json:"id"`
	Fullname     string           `json:"fullname"`
	Photo        string           `json:"photo"`
	MyRecipes    []PersonalRecipe `json:"my_recipes"`
	SavedRecipes []SaveRecipe     `json:"saved_recipes"`
	SavedLikes   []LikeRecipe     `json:"liked_recipes"`
}

type PersonalRecipe struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

type SaveRecipe struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

type LikeRecipe struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}
