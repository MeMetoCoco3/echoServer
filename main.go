package main

import (
	"log"
	"os"
)

func main() {
	store, err := NewBoltStore[string, User](path, bucketName)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	id1, u1, err := NewUser("Vidal", "Maquine", "doneskroto@gmail.com", 12)
	id2, u2, err := NewUser("Katerina", "Principesa", "blabla@gmail.com", 23)
	id3, u3, err := NewUser("Mare", "Reina", "blabla@gmail.com", 34)
	id4, u4, err := NewUser("Pare", "The fucking king", "blabla@gmail.com", 12)
	id5, u5, err := NewUser("Hermano", "Maquine", "blabla@gmail.com", 12)
	id6, u6, err := NewUser("El abuelo Vidal", "Legend of the game", "blabla@gmail.com", 12)
	log.Println(id1)
	log.Println(id2)
	log.Println(id3)
	log.Println(id4)
	log.Println(id5)
	log.Println(id6)
	store.Put(id1.String(), *u1)
	store.Put(id2.String(), *u2)
	store.Put(id3.String(), *u3)
	store.Put(id4.String(), *u4)
	store.Put(id5.String(), *u5)
	store.Put(id6.String(), *u6)

	server, err := NewServerBU(":1337", store)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println(store.Get(id3.String()))

	if err = server.StartServer(); err != nil {
		log.Println("We got some errors: ", err)

	}
}
