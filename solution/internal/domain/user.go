package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"client_id"`
	Login    string    `json:"login"`
	Age      int32     `json:"age"`
	Location string    `json:"location"`
	Gender   string    `json:"gender"`
}
