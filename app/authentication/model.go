package authentication

type SignUpRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Fullname string `json:"fullname" validate:"required"`
	Phone    string `json:"phone"  validate:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInResponse struct {
	Email                   string  `json:"email"`
	AccessToken             string  `json:"access_token"`
	AccessTokenTimeExpired  *string `json:"access_token_time_expired,omitempty"`
	RefreshToken            *string `json:"refresh_token,omitempty"`
	RefreshTokenTimeExpired *string `json:"refresh_token_time_expired,omitempty"`
}
