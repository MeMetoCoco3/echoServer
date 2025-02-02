package main

import (
	"fmt"
	cMiddleware "github.com/MeMetoCoco3/echoServer/middleware"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func (s *ServerBU) handleLoginGet(c echo.Context) error {
	res := &Response{IsLoggedIn: false}
	SetIsLogged(c.Get("user"), res)

	return c.Render(http.StatusOK, "login.html", res)
}

func (s *ServerBU) handleLoginPost(c echo.Context) error {
	res := &Response{IsLoggedIn: false}
	SetIsLogged(c.Get("user"), res)

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
	log.Println("Pre MakeJWT")
	t, err := cMiddleware.MakeJWT(uuid, s.secretKey, time.Duration(time.Minute*2))
	log.Println(err)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]error{"msg": fmt.Errorf("Error generating token: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": t.SignedString})
}

func (s *ServerBU) handleRegister(c echo.Context) error {
	res := &Response{IsLoggedIn: false}
	SetIsLogged(c.Get("user"), res)

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
