package dto

type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
