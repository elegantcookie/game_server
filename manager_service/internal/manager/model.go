package manager

import (
	"context"
	"log"
	"sync"
	"time"
)

type LobbyRecord struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	LobbyID    string `json:"lobby_id" bson:"lobby_id"`
	Expiration int64  `json:"expiration,string" bson:"expiration"`
}

func (lr LobbyRecord) Expired() bool {
	return time.Now().Unix() >= lr.Expiration
}

type LobbyRecordDTO struct {
	LobbyID    string `json:"lobby_id"`
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

//
//func (fq *FuncArray) Pop() *Node {
//	if fq.Head == nil {
//		return nil
//	}
//	head := fq.Head
//	fq.Head = fq.Head.Next
//	fq.Length -= 1
//	return head
//}
