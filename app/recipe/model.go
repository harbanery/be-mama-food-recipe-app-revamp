package recipe

import (
	"time"
)

type Author struct {
	ID       string `json:"id"`
	Fullname string `json:"fullname"`
}

type Save struct {
	RecipeID string `json:"recipe_id"`
	UserID   string `json:"user_id"`
}

type Like struct {
	RecipeID string `json:"recipe_id"`
	UserID   string `json:"user_id"`
}

type Comment struct {
	RecipeID    string `json:"recipe_id"`
	UserID      string `json:"user_id"`
	User        User   `json:"user" gorm:"foreignKey:UserID;"`
	Description string `json:"description"`
}

type Video struct {
	ID       string `json:"id"`
	RecipeID string `json:"recipe_id"`
	URL      string `json:"url"`
}

type VideoRequest struct {
	Title  string `json:"title" validate:"required"`
	Source string `json:"source" validate:"required"`
	URL    string `json:"url" validate:"required"`
}

type VideoResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Source    string    `json:"source"`
	URL       string    `json:"url"`
}

type ActionRecipeResponse struct {
	IsSaved bool `json:"is_saved"`
	IsLiked bool `json:"is_liked"`
}

type User struct {
	ID       string `json:"id"`
	Fullname string `json:"fullname"`
	Photo    string `json:"photo"`
}

type RecipeRequest struct {
	ID          string   `json:"id"`
	Title       string   `json:"title" validate:"required"`
	SubTitle    string   `json:"sub_title"`
	Image       string   `json:"image" validate:"required"`
	Video       []string `json:"video"`
	Header      string   `json:"header"`
	Description string   `json:"description" validate:"required"`
	// Image       []*multipart.FileHeader `json:"image"  validate:"required"`
	// Video       []*multipart.FileHeader `json:"video"`
}

type RecipeResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
	Image  string `json:"image"`
	Header string `json:"header"`
}

type DetailRecipeResponse struct {
	ID          string     `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	SubTitle    string     `json:"sub_title"`
	AuthorID    string     `json:"author_id"`
	Author      *Author    `json:"author"`
	Slug        string     `json:"slug"`
	Header      string     `json:"header"`
	Image       string     `json:"image"`
	Videos      []*Video   `json:"videos" gorm:"foreignKey:RecipeID"`
	Description string     `json:"description"`
	Saves       []*Save    `json:"saves" gorm:"foreignKey:RecipeID"`
	Likes       []*Like    `json:"likes" gorm:"foreignKey:RecipeID"`
	Comments    []*Comment `json:"comments" gorm:"foreignKey:RecipeID"`
}
