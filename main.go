package main

import (
	"log"
	"os"
)

func main() {
	os.Remove(pathUsers)
	os.Remove(pathEmailIndex)

	storeU, err := NewBoltStore[string, User](pathUsers, bucketUsersName)
	storeE, err := NewBoltStore[string, string](pathEmailIndex, bucketEmailName)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	id1, u1, err := NewUser("Vidal", "Maquine", "doneskroto@gmail.com", 12)

	id2, u2, err := NewUser("Katerina", "Principesa", "kata@gmail.com", 23)
	id3, u3, err := NewUser("Mare", "Reina", "Mare@gmail.com", 34)
	id4, u4, err := NewUser("Pare", "The fucking king", "Pare@gmail.com", 12)
	id5, u5, err := NewUser("Hermano", "Maquine", "germa@gmail.com", 12)
	id6, u6, err := NewUser("El abuelo Vidal", "Legend of the game", "AbueloVidal@gmail.com", 12)
	log.Println(id1)
	storeU.Put(id1.String(), *u1)
	storeE.Put(u1.Email, id1.String())
	storeU.Put(id2.String(), *u2)
	storeE.Put(u2.Email, id2.String())
	storeU.Put(id3.String(), *u3)
	storeE.Put(u3.Email, id3.String())
	storeU.Put(id4.String(), *u4)
	storeE.Put(u4.Email, id4.String())
	storeU.Put(id5.String(), *u5)
	storeE.Put(u5.Email, id5.String())
	storeU.Put(id6.String(), *u6)
	storeE.Put(u6.Email, id6.String())
	m, err := storeU.GetValues(id4.String(), []string{"Age", "Email", "Name", "ID", "Role"})
	if err != nil {
		log.Println(err)
	}
	log.Println(m)
	for k, v := range m {
		log.Printf("For K %v, we got %v", k, v)
	}

	/*
	   scan := bufio.NewScanner(os.Stdin)

	   for scan.Scan() {

	   		id, err := storeE.Get(scan.Text())
	   		if err != nil {
	   			log.Println(err)
	   			continue
	   		}

	   		vidalUser, err := storeU.Get(id)
	   		if err != nil {
	   			log.Println(err)
	   			continue
	   		}
	   		log.Println(vidalUser)
	   	}

	   	server, err := NewServerBU(":1337", storeU, storeE)
	   	if err != nil {
	   		log.Println(err)
	   		os.Exit(1)
	   	}
	   	   log.Println(store.Get(id3.String()))

	   	   	if err = server.StartServer(); err != nil {
	   	   		log.Println("We got some errors: ", err)

	   	   }
	*/
}
