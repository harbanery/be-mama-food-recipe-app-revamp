package signup

type SignUp struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Username string `json:"username" validate:"required"`
	Phone    string `json:"phone"  validate:"required"`
}