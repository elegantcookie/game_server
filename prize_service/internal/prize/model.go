package prize

type Prize struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	GameType  string `json:"game_type" bson:"game_type"`
	Result    string `json:"result" bson:"result"`
	TopPlaces string `json:"top_places" bson:"top_places"`
	Reward    string `json:"reward" bson:"reward"`
	DateTime  string `json:"date_time" bson:"date_time"`
}

type PrizeDTO struct {
	GameType  string `json:"game_type" bson:"game_type"`
	Result    string `json:"result" bson:"result"`
	TopPlaces string `json:"top_places" bson:"top_places"`
	Reward    string `json:"reward" bson:"reward"`
	DateTime  string `json:"date_time" bson:"date_time"`
}
