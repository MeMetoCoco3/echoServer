package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *ServerBU) handleLogInGet(c echo.Context) error {
	return c.Render(http.StatusOK, "logIn.html", nil)

}

func (s *ServerBU) handleLogInPost(c echo.Context) error {

	log.Println("Login start")
	// TODO: Bcrypt over here! and on user creation
	email := c.FormValue("email")
	pass := c.FormValue("password")
	log.Println(email)
	uuid, err := s.EmailIndex.Get(email)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"msg": "Email not found!"})
	}

	savedPassword, err := s.Storage.GetValues(uuid, []string{"Password"})
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"msg": "User does not have a password"})
	}

	if pass != savedPassword["Password"] {
		return c.JSON(http.StatusConflict, map[string]string{"msg": "Password does not match"})
	}

	log.Println("Login finish")
	return c.JSON(http.StatusOK, map[string]string{"msg": "Log in good!"})

}
