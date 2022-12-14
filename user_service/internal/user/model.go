package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string        `json:"id" bson:"_id,omitempty"`
	Username      string        `json:"username" bson:"username"`
	Password      string        `json:"-" bson:"password"`
	HasFreeTicket bool          `json:"has_free_ticket" bson:"has_free_ticket"`
	Tickets       []GameTickets `json:"tickets" bson:"tickets"`
}

type TicketDTO struct {
	ID       string `json:"id"`
	GameType string `json:"game_type"`
	TicketID string `json:"ticket_id"`
}

type CreateUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func NewUser(dto CreateUserDTO) User {
	return User{
		Username:      dto.Username,
		Password:      dto.Password,
		HasFreeTicket: true,
		Tickets:       []GameTickets{},
	}
}

//func UpdatedUser(dto UpdateUserDTO) User {
//	return User{
//		ID:       dto.ID,
//		Password: dto.Password,
//	}
//}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match")
	}
	return nil
}

func (u *User) GeneratePasswordHash() error {
	pwd, err := generatePasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password due to error %w", err)
	}
	return string(hash), nil
}
