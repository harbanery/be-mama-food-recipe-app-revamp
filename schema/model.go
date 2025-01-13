package schema

import "mama-recipe/helper"

type User struct {
	helper.UUID
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Username string  `json:"username"`
	Phone    string  `json:"phone"`
	Photo    *string `json:"photo"`
}
