package dto

type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserScoreRequest struct {
	GameID int32 `json:"game_id" validate:"required"`
}

type UserScoreResponse struct {
	UserID int32 `json:"user_id"`
	GameID int32 `json:"game_id"`
	Score  int32 `json:"score"`
}
