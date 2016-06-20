package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

//Import the json configuration file.
func importConfiguration(configPath string) (configuration, error) {
	fmt.Printf("Importing the configuration information from %v\n", configPath)

	f, err := ioutil.ReadFile(configPath)
	var c configuration

	if err != nil {
		fmt.Printf("Done. Error %s.\n", err.Error())
		return c, err
	}

	json.Unmarshal(f, &c)

	fmt.Printf("Done. Configuration data: \n %+v \n", c)

	return c, err
}

//Returns "in" or "out" depending on the current status of the user
func getStatus(userID string) (string, error) {
	fmt.Printf("Getting Status of user %s\n", userID)
	usr, err := getUserByID(userID)

	if err != nil {
		fmt.Printf("Done. Error %s.\n", err.Error())
		return "", err
	}
	if strings.EqualFold(usr.Status, "") {
		//Just for completions sake, since we should always have a status as users
		//are created with status
		err := errors.New("Could not retrieve user Status, user has no status.")
		fmt.Printf("Done. Error %s.\n", err.Error())
		return "", err
	}

	fmt.Printf("Done. Status %s.\n", userID)
	return usr.Status, nil
}

func clockOut(usr user) error {
	fmt.Printf("Clocking user Out %s.\n", usr.ID)
	status, err := getStatus(usr.ID)
	if err != nil {
		fmt.Printf("Done. Error %s.\n", err.Error())
		return err
	}

	//Can't clock in if we're already in.
	//
	//TODO: Think about this? Do we want to allow users to clock in anyway and
	//create a missed punch like thing?
	if strings.EqualFold(status, "out") {
		err := errors.New("Could not clock in, already clocked in.")
		fmt.Printf("Done. Error %s.\n", err.Error())
		return err
	}
	//Do other checking?
	return setUserStatus(usr.ID, "out")
}

//Return nil if successful, else returns an error.
func clockIn(usr user) error {
	fmt.Printf("Clocking user in %s.\n", usr.ID)
	status, err := getStatus(usr.ID)
	if err != nil {
		fmt.Printf("Done. Error %s.\n", err.Error())
		return err
	}

	//Can't clock in if we're already in.
	//
	//TODO: Think about this? Do we want to allow users to clock in anyway and
	//create a missed punch like thing?
	if strings.EqualFold(status, "in") {
		err := errors.New("Could not clock out, already clocked out.")
		fmt.Printf("Done. Error %s.\n", err.Error())
		return err
	}
	//Do other checking?
	return setUserStatus(usr.ID, "in")
}

//check to make sure the required fields of a user are present.
func validateUser(u *user) error {
	if strings.EqualFold(u.Name, "") {
		return errors.New("No name provided.")
	}
	return nil
}
