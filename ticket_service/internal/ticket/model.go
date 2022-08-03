package ticket

type Ticket struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	IsActive     bool   `json:"is_active" bson:"is_active"`
	TicketPrice  int    `json:"ticket_price" bson:"ticket_price"`
	PlayerAmount int    `json:"player_amount" bson:"player_amount"`
	GameType     string `json:"game_type" bson:"game_type"`
	PrizeId      string `json:"prize_id" bson:"prize_id"`
}

type TicketDTO struct {
	TicketPrice  int    `json:"ticket_price"`
	PlayerAmount int    `json:"player_amount"`
	GameType     string `json:"game_type"`
	PrizeId      string `json:"prize_id"`
}

type TicketIDDTO struct {
	ID string `json:"id"`
}

type FreeTicketStatusDTO struct {
	Status    bool   `json:"tickets_available"`
	AccessKey string `json:"access_key"`
}
