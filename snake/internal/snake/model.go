package snake

type Snake struct {
	ID        string   `json:"id" bson:"_id,omitempty"`
	Players   []string `json:"players" bson:"players"`
	Results   []Player `json:"results" bson:"results"`
	StartTime int64    `json:"start_time" bson:"start_time"`
	EndTime   int64    `json:"end_time" bson:"end_time"`
}

func NewSnake(dto SnakeDTO) Snake {
	return Snake{
		Players:   dto.Players,
		Results:   nil,
		StartTime: dto.StartTime,
		EndTime:   dto.EndTime,
	}
}

type Player struct {
	UserID string `json:"user_id"`
	Result int    `json:"result"`
}

func NewPlayer(dto SendResultDTO) Player {
	return Player{
		UserID: dto.UserID,
		Result: dto.Result,
	}
}

type SnakeDTO struct {
	Players   []string `json:"players" bson:"players"`
	StartTime int64    `json:"start_time" bson:"start_time"`
	EndTime   int64    `json:"end_time" bson:"end_time"`
}

type SendResultDTO struct {
	GameServerID string `json:"id"`
	UserID       string `json:"user_id"`
	Result       int    `json:"result"`
}
