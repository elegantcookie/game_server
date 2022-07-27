package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"-" bson:"password"`
}

type CreateUserDTO struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password" bson:"-"`
}

type UpdateUserDTO struct {
	ID          string `json:"id,omitempty" bson:"_id,omitempty"`
	Password    string `json:"password,omitempty" bson:"password,omitempty"`
	OldPassword string `json:"old_password,omitempty" bson:"-"`
	NewPassword string `json:"new_password,omitempty" bson:"-"`
}

func NewUser(dto CreateUserDTO) User {
	return User{
		Username: dto.Username,
		Password: dto.Password,
	}
}

func UpdatedUser(dto UpdateUserDTO) User {
	return User{
		ID:       dto.ID,
		Password: dto.Password,
	}
}

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
