package main

import (
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *ServerBU) handlePut(c echo.Context) error {
	name := c.Param("name")
	role := c.Param("role")

	age, err := strconv.Atoi(c.Param("age"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": err})

	}
	// TODO
	defaultEmail := "robocop@gmail.com"
	uuid, newUser, err := NewUser(name, role, defaultEmail, age)

	err = c.Validate(newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]error{"msg": err})
	}

	s.Storage.Put(uuid.String(), *newUser)
	msg := fmt.Sprintf("Transaction completed: Added User{'name':'%v','age':%v,'role':'%v'}", name, age, role)
	return c.JSON(http.StatusOK, map[string]string{"msg": msg})
}

func (s *ServerBU) handleGet(c echo.Context) error {
	uuid := c.Param("id")
	if uuid == "" {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": fmt.Errorf("Id was not introduced")})
	}

	user, err := s.Storage.Get(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": err})
	}

	//userBytes, err := json.Marshal(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": err})
	}

	templateName := "GetUser.html"
	log.Println("Rendering template:, ", templateName)
	if err = c.Render(http.StatusOK, templateName, user); err != nil {
		log.Println(err)
	}
	return err
	//return c.JSON(http.StatusOK, map[string]interface{}{"user": json.RawMessage(userBytes)})
}

func (s *ServerBU) handleGetAll(c echo.Context) error {
	users, err := s.Storage.GetAll()
	if err != nil {
		log.Println("Error fetching all users:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "Error fetching users"})
	}
	templateName := "GetAllUsers.html"
	log.Println("Rendering template:, ", templateName)
	if err = c.Render(http.StatusOK, templateName, users); err != nil {
		log.Println(err)
	}
	return err
}

func (s *ServerBU) handleDelete(c echo.Context) error {
	uuid := c.Param("id")
	if uuid == "" {
		uuid = c.FormValue("id")
	}

	if uuid == "" {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": fmt.Errorf("Id was not introduced")})
	}

	err := s.Storage.Delete(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": fmt.Errorf("UUID: %v, Error: %v", uuid, err)})
	}

	msg := fmt.Sprintf("User with id %v was deleted from the Database.", uuid)
	return c.JSON(http.StatusOK, map[string]string{"msg": msg})
}

func (s *ServerBU) handleUpdateUserData(c echo.Context) error {
	field := c.Param("field")
	id := c.Param("id")
	var info struct {
		Value string `json:"value"`
	}
	if err := c.Bind(&info); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]error{"err": err})
	}

	user, err := s.Storage.Get(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]error{"err": err})
	}
	switch field {
	case "role":
		user.Role = info.Value
	case "age":
		ageVal, err := strconv.Atoi(info.Value)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]error{"err": err})
		}
		user.Age = ageVal
	case "email":
		user.Email = info.Value
	case "description":
		user.Description = info.Value
	}

	err = s.Storage.Put(id, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]error{"err": err})
	}
	return c.JSON(http.StatusOK, map[string]bool{"success": true})

}
