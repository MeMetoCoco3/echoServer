package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	cMiddleware "github.com/MeMetoCoco3/echoServer/middleware"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	_ "github.com/labstack/echo-jwt/v4"
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

type (
	ServerBU struct {
		laddr      string
		Storage    Storer[string, User]   // uuid to User
		EmailIndex Storer[string, string] // email to uuid
		secretKey  string
	}

	CustomValidator struct {
		Validator *validator.Validate
	}

	Template struct {
		templates *template.Template
	}
)

func NewServerBU(laddr string, store Storer[string, User], emailIndex Storer[string, string]) (*ServerBU, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("No key was found as envVar")
	}

	return &ServerBU{
		laddr:      laddr,
		Storage:    store,
		EmailIndex: emailIndex,
		secretKey:  secretKey,
	}, nil
}

func (s *ServerBU) StartServer() error {
	e := echo.New()

	templates := template.Must(template.ParseGlob("templates/*.html"))
	t := &Template{
		templates: templates,
	}
	e.Renderer = t

	e.Validator = CustomValidator{Validator: validator.New()}
	e.Use(cMiddleware.ResponseLogger)
	e.Use(cMiddleware.RealIP)
	e.Use(cMiddleware.OptionalJWT([]byte(s.secretKey)))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:1337"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	//e.Use(middleware.LoggerWithConfig(CustomLoggerConfig))

	e.Static("/static", "static")
	e.GET("/about", func(c echo.Context) error {
		return c.Render(http.StatusOK, "about.html", nil)
	})
	e.GET("/get/:id", s.handleGet)
	e.GET("/getAll", s.handleGetAll)
	e.GET("/register", s.handleLoginGet)
	e.POST("/register", s.handleRegister)
	e.POST("/login", s.handleLoginPost)
	e.PUT("/put/:name/:role/:age", s.handlePut)
	e.POST("/update/:id/:field", s.handleUpdateUserData)
	e.POST("/delete/:id", s.handleDelete)

	return e.Start(s.laddr)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if err := t.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (cv CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Not enought required data"))
	}
	return nil
}

func CookieHeaders(value, name string) func(echo.Context) error {
	writeCookie := func(c echo.Context) error {
		cookie := new(http.Cookie)
		cookie.Value = value
		cookie.Name = name
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
