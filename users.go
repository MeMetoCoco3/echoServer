package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"reflect"
	"strconv"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Role        string    `json:"role" validate:"required"`
	Age         int       `json:"age" validate:"required"`
	Email       string    `json:"email" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Password    string    `json:"password"`
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
		Password:    "1234",
	}, nil
}

// ----------------------------- This "u any" saved my life
// --------------------------------------I------------------
// --------------------------------------I------------------
// -------------------------------------\I/-----------------
// --------------------------------------V------------------
func FilterStruct[K comparable, V any](u any, keys []K) map[K]string {
	log.Println("Start filtering")
	result := make(map[K]string)
	v := reflect.ValueOf(u)

	// Checks if pointer, if yes, dereference
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Checks if is a struct
	if v.Kind() != reflect.Struct {
		return result
	}

	for _, field := range keys {
		fieldName, ok := any(field).(string)
		if !ok {
			continue
		}

		f := v.FieldByName(fieldName)
		if !f.IsValid() {
			continue
		}

		var value string
		//check if is UUID
		log.Printf("f.Type().String(): %v\n", f.Type().String())
		if f.Type().String() == "uuid.UUID" {
			if uuid, ok := f.Interface().(uuid.UUID); ok {
				result[field] = uuid.String()
				continue
			}
		}

		log.Printf("Value of field %v = %v\n", fieldName, f.String())
		switch f.Kind() {
		case reflect.String:
			value = f.String()
		case reflect.Int:
			value = strconv.Itoa(int(f.Int()))
		default:
			log.Println("continue")
			continue
		}

		result[field] = value
	}

	return result
}

func (u *User) JSON() ([]byte, error) {
	return json.Marshal(u)
}
