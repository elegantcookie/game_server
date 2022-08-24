package manager

import (
	"context"
	"log"
	"sync"
	"time"
)

type LobbyRecord struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Type       string `json:"type" bson:"type"`
	LobbyID    string `json:"lobby_id" bson:"lobby_id"`
	GameType   string `json:"game_type" bson:"game_type"`
	Expiration int64  `json:"expiration,string" bson:"expiration"`
}

func (lr LobbyRecord) Expired() bool {
	return time.Now().Unix() >= lr.Expiration
}

type LobbyRecordDTO struct {
	Type       string `json:"type"`
	LobbyID    string `json:"lobby_id"`
	GameType   string `json:"game_type"`
	Expiration int64  `json:"expiration"`
	JWTToken   string `json:"-"`
}

type UpdateTimeDTO struct {
	Expiration int64  `json:"expiration"`
	JWTToken   string `json:"-"`
}

type FuncArray struct {
	ManagerService Service
	Length         int
	*sync.Mutex
}

var instance FuncArray
var once sync.Once

func GetFuncQueue(managerService Service) FuncArray {
	once.Do(func() {
		instance = FuncArray{
			ManagerService: managerService,
			Length:         0,
		}
	})

	return instance
}

func (fq *FuncArray) Update(lr LobbyRecord) error {
	log.Println("UPDATE")
	err := fq.ManagerService.UpdateLR(context.Background(), lr)
	if err != nil {
		return err
	}
	return nil
}

type LRResponse struct {
	UpdatedTime int64
	StatusCode  int
	Delete      bool
}

func (r LRResponse) CorrectResponse() bool {
	if r.StatusCode > 299 || r.StatusCode == 0 {
		return false
	}
	return true
}
