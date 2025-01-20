package schema

import "mama-recipe/helper"

type User struct {
	helper.UUID
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Username string    `json:"username"`
	Fullname string    `json:"fullname"`
	Phone    string    `json:"phone"`
	Photo    string    `json:"photo"`
	Recipes  []*Recipe `json:"recipes" gorm:"foreignKey:AuthorID"`
	Saves    []*Save   `json:"saves" gorm:"foreignKey:UserID"`
	Likes    []*Like   `json:"likes" gorm:"foreignKey:UserID"`
}

type Recipe struct {
	helper.UUID
	Title       string     `json:"title"`
	SubTitle    string     `json:"sub_title"`
	Slug        string     `json:"slug"`
	Header      string     `json:"header"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	Videos      []*Video   `json:"videos" gorm:"foreignKey:RecipeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Author      User       `json:"author" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	AuthorID    string     `json:"author_id"`
	Saves       []*Save    `json:"saves" gorm:"foreignKey:RecipeID"`
	Likes       []*Like    `json:"likes" gorm:"foreignKey:RecipeID"`
	Comments    []*Comment `json:"comments" gorm:"foreignKey:RecipeID"`
}

type Save struct {
	helper.UUID
	RecipeID string `json:"recipe_id"`
	Recipe   Recipe `json:"recipe" gorm:"foreignKey:RecipeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID   string `json:"user_id"`
	User     User   `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Like struct {
	helper.UUID
	RecipeID string `json:"recipe_id"`
	Recipe   Recipe `json:"recipe" gorm:"foreignKey:RecipeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID   string `json:"user_id"`
	User     User   `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Comment struct {
	helper.UUID
	RecipeID    string `json:"recipe_id"`
	Recipe      Recipe `json:"recipe" gorm:"foreignKey:RecipeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID      string `json:"user_id"`
	User        User   `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Description string `json:"description"`
}

type Video struct {
	helper.UUID
	RecipeID string `json:"recipe_id"`
	Recipe   Recipe `json:"recipe" gorm:"foreignKey:RecipeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title    string `json:"title"`
	Source   string `json:"source"`
	URL      string `json:"url"`
}
