package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
)

func clockInHandler(c echo.Context) error {
	usr := user{ID: c.Param("userID")}

	err := clockIn(usr)
	if err == nil {
		c.Response().Write([]byte("Clocked in."))
	}
	return err
}

func clockOutHandler(c echo.Context) error {
	usr := user{ID: c.Param("userID")}

	err := clockOut(usr)
	if err == nil {
		c.Response().Write([]byte("Clocked out."))
	}
	return err
}

func postUserHandler(c echo.Context) error {
	fmt.Printf("Posting a user\n")
	b, _ := ioutil.ReadAll(c.Request().Body())

	var u user

	json.Unmarshal(b, &u)
	err := validateUser(&u)
	if err != nil {
		return err
	}

	u.Status = "Out"

	return createUser(&u)
}
