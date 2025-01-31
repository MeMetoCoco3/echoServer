package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *ServerBU) handleLogInGet(c echo.Context) error {
	return c.Render(http.StatusOK, "logIn.html", nil)

}

func (s *ServerBU) handleLogInPost(c echo.Context) error {
	email := c.FormValue("email")
	pass := c.FormValue("password")
	uuid, err := s.EmailIndex.Get(email)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"msg": "Email not found!"})
	}

	savedPassword, err := s.Storage.GetValues(uuid, []string{"Password"})
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"msg": "User does not have a password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedPassword["Password"]), []byte(pass))
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"msg": "Password does not match"})
	}

	return c.JSON(http.StatusOK, map[string]string{"msg": "Log in good!"})
}

func (s *ServerBU) handleRegister(c echo.Context) error {
	// TODO: We expect every value is valid
	email := c.FormValue("email")
	pass := c.FormValue("password")
	username := c.FormValue("username")

	uuid, u, err := NewUser(username, "default", email, pass, 0)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"msg": "Error creating user"})
	}

	err = s.Storage.Put(uuid.String(), *u)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"msg": "Error putting user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"msg": fmt.Sprintf("User %v was registered!!", username)})
}
