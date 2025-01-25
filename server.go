package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const boltStoreName = "bunny"

// Logger middleware configuration

var formatVerbose = `{"time":"${time_rfc3339_nano}", "status":${status}, "remote_ip":"${remote_ip}", ` +
	`"method":"${method}", "host":"${host}", "uri":"${uri}", "user_agent":"${user_agent}",` +
	`"latency":${latency}, "latency_human":"${latency_human}",` +
	`"bytes_in":${bytes_in}, "bytes_out":${bytes_out}, "error":"${error}"}` + "\n"

var format = `${method} ${status} ${uri}. Error: "${error}"` + "\n"

var CustomLoggerConfig = middleware.LoggerConfig{
	Skipper:          middleware.DefaultSkipper,
	Format:           format,
	CustomTimeFormat: "2006-01-02 15:04:05.00000",
}

type ServerBU struct {
	laddr   string
	Storage Storer[string, User]
}

func NewServerBU(laddr string, store Storer[string, User]) (*ServerBU, error) {
	return &ServerBU{
		laddr:   laddr,
		Storage: store,
	}, nil
}

func (s *ServerBU) StartServer() error {

	log.Println("Starting server on: %s.", s.laddr)

	e := echo.New()

	e.Static("/static", "static")

	e.Use(middleware.LoggerWithConfig(CustomLoggerConfig))
	e.Use(RealIPMiddleware)
	e.PUT("/put/:name/:role/:age", s.handlePut)
	e.GET("/get/:id", s.handleGet)
	e.GET("/getAll", s.handleGetAll)
	e.POST("/delete/:id", s.handleDelete)

	return e.Start(s.laddr)
}

func CookieHeaders(values map[string]string) func(echo.Context) error {
	writeCookie := func(c echo.Context) error {
		cookie := new(http.Cookie)
		cookie.Name = values["name"]
		cookie.Value = values["pass"]
		cookie.Expires = time.Now().Add(20 * time.Second)

		c.SetCookie(cookie)
		return c.String(http.StatusOK, "write a cookie")
	}
	return writeCookie
}

/*
	type Server[K comparable, V any] struct {
		laddr   string
		Storage Storer[K, V]
	}
*/
