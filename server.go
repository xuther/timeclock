package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

var config = configuration{}

func main() {

	c, err := importConfiguration("./config.json")
	if err != nil {
		fmt.Printf("Could not import the configuration file, check that it exists: %s\n", err.Error())
	}
	config = c

	e := echo.New()

	e.Post("/api/clockin", clockIn)
	e.Post("/api/users", postUser)
	e.Static("/pages", "Static")
	e.Static("/scripts", "Static/scripts")

	e.Run(standard.New(":8888"))

	//eventually we need to accept command line parameters.
}
