package recipe

type Author struct {
	ID       string `json:"id"`
	Fullname string `json:"fullname"`
}

type RecipeRequest struct {
	Title       string `json:"title" validate:"required"`
	SubTitle    string `json:"sub_title" validate:"required"`
	Image       string `json:"image"`
	Header      string `json:"header" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type RecipeResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
	Image  string `json:"image"`
	Author string `json:"author"`
}

type AlternativeRecipeResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
	Header string `json:"header"`
	Image  string `json:"image"`
	Author string `json:"author"`
}

type DetailRecipeResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	SubTitle    string `json:"sub_title"`
	Author      string `json:"author"`
	Slug        string `json:"slug"`
	Header      string `json:"header"`
	Image       string `json:"image"`
	Description string `json:"description"`
}
