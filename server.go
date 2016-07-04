package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

var config = configuration{}

func main() {

	c, err := importConfiguration("./config.json")
	if err != nil {
		fmt.Printf("Could not import the configuration file, check that it exists: %s\n", err.Error())
	}
	config = c

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.Post("/api/users", postUserHandler)
	e.Post("/api/users/:userID/clockin", clockInHandler)
	e.Post("/api/users/:userID/clockout", clockOutHandler)

	e.Get("/api/users/:userID", getUserHandler)
	e.Get("/api/users/:userID/lastpunch", getLastPunchHandler)
	e.Get("/api/users/:userID/punches", getPunchesHandler)

	e.Static("/pages", "Static")
	e.Static("/scripts", "Static/scripts")

	e.Run(standard.New(":8888"))

	//eventually we need to accept command line parameters.
}
