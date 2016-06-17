package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
)

func clockIn(c echo.Context) error {
	b, _ := ioutil.ReadAll(c.Request().Body())
	fmt.Printf("%s", b)

	return nil
}

func postUser(c echo.Context) error {
	b, _ := ioutil.ReadAll(c.Request().Body())

	var u user

	json.Unmarshal(b, &u)

	return createUser(&u)
}
