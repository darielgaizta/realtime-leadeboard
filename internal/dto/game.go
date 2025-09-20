package dto

import db "github.com/darielgaizta/realtime-leaderboard/internal/db/generated"

type GameRequest struct {
	Name string `json:"name" validate:"required"`
}

type GameResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func ToGameResponses(games []db.Game) []GameResponse {
	responses := make([]GameResponse, 0, len(games))
	for _, game := range games {
		responses = append(responses, GameResponse{
			ID:   game.ID,
			Name: game.Name,
		})
	}
	return responses
}

type UserGameScoreRequest struct {
	Score int32 `json:"score"`
}

type UserGameScoreResponse struct {
	UserID int32 `json:"user_id"`
	GameID int32 `json:"game_id"`
	Score  int32 `json:"score"`
}
