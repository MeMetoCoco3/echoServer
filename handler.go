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
	res := &Response{IsLoggedIn: false}
	SetIsLogged(c.Get("user"), res)

	name := c.Param("name")
	role := c.Param("role")
	password := c.Param("password")
	//email := c.Param("email")
	defaultEmail := "robocop@gmail.com"

	age, err := strconv.Atoi(c.Param("age"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": err})

	}
	uuid, newUser, err := NewUser(name, role, defaultEmail, password, age)

	err = c.Validate(newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]error{"msg": err})
	}

	s.Storage.Put(uuid.String(), *newUser)
	s.EmailIndex.Put(defaultEmail, uuid.String())
	msg := fmt.Sprintf("Transaction completed: Added User{'name':'%v','age':%v,'role':'%v', 'password': '%v'}", name, age, role, password)
	return c.JSON(http.StatusOK, map[string]string{"msg": msg})
}

func (s *ServerBU) handleGet(c echo.Context) error {
	res := &Response{IsLoggedIn: false}
	SetIsLogged(c.Get("user"), res)

	uuid := c.Param("id")
	if uuid == "" {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": fmt.Errorf("Id was not introduced")})
	}

	user, err := s.Storage.Get(uuid)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]error{"msg": fmt.Errorf("UUid not found: %v", err)})
	}

	res.Content = user
	return c.Render(http.StatusOK, "GetUser.html", res)
}

func (s *ServerBU) handleGetAll(c echo.Context) error {
	res := &Response{IsLoggedIn: false}
	SetIsLogged(c.Get("user"), res)

	log.Println("C get is equal to: ", c.Get("user"))

	users, err := s.Storage.GetAll()
	if err != nil {
		log.Println("Error fetching all users:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "Error fetching users"})
	}
	res.Content = users
	return c.Render(http.StatusOK, "GetAllUsers.html", res)
}

func (s *ServerBU) handleDelete(c echo.Context) error {
	res := &Response{IsLoggedIn: false}
	SetIsLogged(c.Get("user"), res)

	uuid := c.Param("id")
	if uuid == "" {
		uuid = c.FormValue("id")
	}

	if uuid == "" {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": fmt.Errorf("Id was not introduced")})
	}
	userInfo, err := s.Storage.Get(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": fmt.Errorf("Id was not introduced")})
	}
	err = s.EmailIndex.Delete(userInfo.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": fmt.Errorf("Id was not introduced")})
	}
	err = s.Storage.Delete(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]error{"msg": fmt.Errorf("UUID: %v, Error: %v", uuid, err)})
	}

	msg := fmt.Sprintf("User with id %v was deleted from the Database.", uuid)
	return c.JSON(http.StatusOK, map[string]string{"msg": msg})
}

func (s *ServerBU) handleUpdateUserData(c echo.Context) error {
	res := &Response{IsLoggedIn: false}
	SetIsLogged(c.Get("user"), res)

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
		if userData, err := s.Storage.Get(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]error{"err": err})
		} else {
			user.Email = info.Value
			err = s.EmailIndex.Delete(userData.Email)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]error{"err": err})
			}
		}
	case "description":
		user.Description = info.Value
	}

	err = s.Storage.Put(id, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]error{"err": err})
	}

	err = s.EmailIndex.Put(user.Email, id)

	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}
