package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {

	log.Println(7 / 2)

	log.Println("start")
	password := []byte("UnPassRandom123")
	password2 := []byte("OtroPassjj")

	// Cost factor of 12 is enough, if it is more it can take far longer.
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		log.Println(err)
	}

	hashedPassword2, err := bcrypt.GenerateFromPassword(password2, 12)
	if err != nil {
		log.Println(err)
	}

	log.Printf("We are comparing %v and %v\n", password, hashedPassword)

	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err == nil {
		log.Println("It is good!!")
	} else {
		log.Println("It is not good!")
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword2, password)
	if err == nil {
		log.Println("It is good!!")
	} else {
		log.Println("It is not good!")
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, password2)
	if err == nil {
		log.Println("It is good!!")
	} else {
		log.Println("It is not good!")
	}
}
