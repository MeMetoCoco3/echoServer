package main

import (
	"bytes"
	"encoding/json"
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

	uuid, newUser, err := NewUser(name, role, age)

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

	userBytes, err := json.Marshal(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": err})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"user": json.RawMessage(userBytes)})
}

func (s *ServerBU) handleGetAll(c echo.Context) error {
	users, err := s.Storage.GetAll()
	if err != nil {
		log.Println("Error fetching all users:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "Error fetching users"})
	}

	// DEALING WITH HTML
	var buff bytes.Buffer

	err = IssueList.Execute(&buff, users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": err})
	}
	return c.JSON(http.StatusOK, users)
}

func (s *ServerBU) handleDelete(c echo.Context) error {
	uuid := c.Param("id")
	if uuid != "" {

		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": fmt.Errorf("Id was not introduced")})
	}

	err := s.Storage.Delete(uuid)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": err})
	}

	msg := fmt.Sprintf("User with id %v was deleted from the Database.", uuid)
	return c.JSON(http.StatusOK, map[string]string{"msg": msg})
}
