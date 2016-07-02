package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/labstack/echo"
)

func clockInHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return err
	}

	usr := user{ID: id}

	err = clockIn(usr)
	if err == nil {
		c.Response().Write([]byte("Clocked in."))
	}
	return err
}

func clockOutHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return err
	}

	usr := user{ID: id}

	err = clockOut(usr)
	if err == nil {
		c.Response().Write([]byte("Clocked out."))
	}
	return err
}

func getLastPunchHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return err
	}

	punch, err := getLastTimepunch(id)
	if err != nil {
		return err
	}

	b, err := json.Marshal(&punch)
	if err != nil {
		return err
	}

	c.Response().Write(b)

	return nil
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

	u.Status = false

	return createUser(&u)
}
