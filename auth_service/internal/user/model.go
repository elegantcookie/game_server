package user

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"-" bson:"password"`
}

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseUserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
