package table

type Record struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	UserID    string `json:"user_id" bson:"user_id"`
	Username  string `json:"username" bson:"username"`
	UserScore int    `json:"user_score" bson:"user_score"`
}

func NewRecord(dto RecordDTO) Record {
	return Record{
		UserID:    dto.UserID,
		Username:  dto.Username,
		UserScore: dto.UserScore,
	}
}

type RecordDTO struct {
	TableName string `json:"table_name" bson:"table_name"`
	ID        string `json:"id" bson:"_id,omitempty"`
	UserID    string `json:"user_id" bson:"user_id"`
	Username  string `json:"username" bson:"username"`
	UserScore int    `json:"user_score" bson:"user_score"`
	JWTToken  string `json:"-" bson:"-"`
}

type Collection struct {
	Name string `json:"table_name" bson:"table_name"`
}

type CollectionDTO struct {
	AccessKey string `json:"access_key"`
	Name      string `json:"table_name"`
	JWTToken  string `json:"-"`
}

type CreateTicketDTO struct {
	UserID   string `json:"user_id"`
	GameType string `json:"game_type"`
	JWT      string `json:"-"`
}

// CTicketDTO includes fields to create ticket with changed signature
type CTicketDTO struct {
	UserID       string `json:"user_id"`
	IsGift       bool   `json:"is_gift"`
	TicketPrice  int    `json:"ticket_price"`
	PlayerAmount int    `json:"player_amount"`
	GameType     string `json:"game_type"`
	PrizeId      string `json:"prize_id"`
}

func NewCreateTicketDTO(gameType string, dto RecordDTO) CreateTicketDTO {
	return CreateTicketDTO{
		UserID:   dto.UserID,
		GameType: gameType,
		JWT:      dto.JWTToken,
	}
}

type AddTicketToUserDTO struct {
	ID       string `json:"id"`
	TicketID string `json:"ticket_id"`
	GameType string `json:"game_type"`
	JWTToken string ` json:"-"`
}

type TicketDTO struct {
	TicketID string `json:"ticket_id"`
	JWTToken string ` json:"-"`
}

type NotifyManagerDTO struct {
	Type       string `json:"type"`
	GameType   string `json:"game_type"`
	Expiration int64  `json:"expiration"`
}

func ReverseArray(array []Record) {
	usersLen := len(array)
	for i := 0; i < usersLen/2; i++ {
		array[i], array[usersLen-1-i] = array[usersLen-1-i], array[i]
	}
}
