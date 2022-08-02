package ticket

type Ticket struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Type         string `json:"type" bson:"type"`
	TicketPrice  int    `json:"ticket_price" bson:"ticket_price"`
	PlayerAmount int    `json:"player_amount" bson:"player_amount"`
	GameType     string `json:"game_type" bson:"game_type"`
	PrizeId      string `json:"prize_id" bson:"prize_id"`
}

type TicketDTO struct {
	Type         string `json:"type"`
	TicketPrice  int    `json:"ticket_price"`
	PlayerAmount int    `json:"player_amount"`
	GameType     string `json:"game_type"`
	PrizeId      string `json:"prize_id"`
}

type TicketIDDTO struct {
	ID string `json:"id"`
}
