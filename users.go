package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Role        string    `json:"role" validate:"required"`
	Age         int       `json:"age" validate:"required"`
	Email       string    `json:"email" validate:"required"`
	Description string    `json:"description" validate:"required"`
}

func NewUser(name, role, email string, age int) (uuid.UUID, *User, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return uuid, nil, fmt.Errorf("Could not create UUID: %v", err)
	}
	return uuid, &User{
		ID:          uuid,
		Name:        name,
		Role:        role,
		Age:         age,
		Email:       email,
		Description: "No description... YET!!",
	}, nil
}

func (u *User) JSON() ([]byte, error) {
	bytes, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return bytes, nil

}
