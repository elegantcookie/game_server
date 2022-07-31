package table

type Record struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	UserID    string `json:"user_id" bson:"user_id"`
	UserScore string `json:"user_score" bson:"user_score"`
}

type RecordDTO struct {
	TableName string `json:"table_name" bson:"table_name"`
	ID        string `json:"id" bson:"_id,omitempty"`
	UserID    string `json:"user_id" bson:"user_id"`
	UserScore string `json:"user_score" bson:"user_score"`
}

type Collection struct {
	Name string `json:"table_name" bson:"table_name"`
}

type CollectionDTO struct {
	AccessKey string `json:"access_key"`
	Name      string `json:"table_name"`
}

func NewRecord(dto RecordDTO) Record {
	return Record{
		UserID:    dto.UserID,
		UserScore: dto.UserScore,
	}
}
