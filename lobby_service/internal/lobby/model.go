package lobby

type Lobby struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	GameType    string   `json:"game_type" bson:"game_type"`
	MaxPlayers  int      `json:"max_players" bson:"max_players"`
	NowPlayers  int      `json:"now_players" bson:"now_players"`
	TicketPrice int      `json:"ticket_price" bson:"ticket_price"`
	PrizeSum    int      `json:"prize_sum" bson:"prize_sum"`
	PrizeType   int      `json:"prize_type" bson:"prize_type"`
	Players     []Player `json:"players" bson:"players"`
	StartTime   int64    `json:"start_time" bson:"start_time"`
	EndTime     int64    `json:"end_time" bson:"end_time"`
}

type Player struct {
	ID    string `json:"user_id"`
	Ready bool   `json:"ready"`
}

type UpdateUserDTO struct {
	ID            string        `json:"id" bson:"_id,omitempty"`
	Username      string        `json:"username" bson:"username"`
	HasFreeTicket bool          `json:"has_free_ticket" bson:"has_free_ticket"`
	Tickets       []GameTickets `json:"tickets" bson:"tickets"`
}

type GameTickets struct {
	GameType string   `json:"game_type"`
	Amount   int      `json:"amount"`
	IDsOfGT  []string `json:"tickets_of_gt"`
}

//type CreateLobbyDTO struct {
//	GameType    string `json:"game_type"`
//	MaxPlayers  int    `json:"max_players"`
//	NowPlayers  int    `json:"now_players"`
//	TicketPrice int    `json:"ticket_price"`
//	PrizeSum    int    `json:"prize_sum"`
//	PrizeType   int    `json:"prize_type"`
//}

type LobbyDTO struct {
	GameType    string `json:"game_type"`
	MaxPlayers  int    `json:"max_players"`
	NowPlayers  int    `json:"now_players"`
	TicketPrice int    `json:"ticket_price"`
	PrizeSum    int    `json:"prize_sum"`
	PrizeType   int    `json:"prize_type"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	JWTToken    string `json:"-"`
}

type JoinLobbyDTO struct {
	UserID   string `json:"user_id"`
	LobbyID  string `json:"lobby_id"`
	TicketID string `json:"ticket_id"`
	JWTToken string `json:"-"`
}

type Params struct {
	GameType   string `json:"game_type"`
	PrizeSum   int    `json:"prize_sum"`
	MaxPlayers int    `json:"max_players"`
}

type UpdateTimeDTO struct {
	ID       string `json:"id"`
	JWTToken string `json:"-"`
}
