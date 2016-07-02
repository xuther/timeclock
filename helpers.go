package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
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
func getStatus(userID int64) (bool, error) {
	fmt.Printf("Getting Status of user %v\n", userID)
	usr, err := getUserByID(userID)

	if err != nil {
		fmt.Printf("Done. Error %s.\n", err.Error())
		return false, err
	}

	fmt.Printf("Done. Status %v.\n", userID)
	return usr.Status, nil
}

func clockOut(usr user) error {
	fmt.Printf("Clocking user Out %s.\n", usr.ID)

	//Do other checking?
	punch, err := getLastTimepunch(usr.ID)
	if err != nil {
		return err
	}

	//If the last punch exists, has an in, but not an out,
	//complete the punch.
	if (!punch.In.Equal(time.Time{})) && (punch.Out.Equal(time.Time{})) {
		punch.Out = time.Now()

		err = updatePunch(punch)

		if err != nil {
			return err
		}
	} else { //in every other case, we just want to create a new punch.
		err = createPunch(timePunch{UID: usr.ID, Out: time.Now()})

		if err != nil {
			return err
		}
	}
	fmt.Printf("Done.")
	return setUserStatus(usr.ID, false)
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
	//TODO: Think about this? Do we want to allow users to clock in anyway and
	//create a missed punch like thing?
	if status {
		err := errors.New("Could not clock in, already clocked in.")
		fmt.Printf("Done. Error %s.\n", err.Error())
		return err
	}

	err = createPunch(timePunch{UID: usr.ID, In: time.Now()})
	if err != nil {
		return err
	}

	//Do other checking?
	err = setUserStatus(usr.ID, true)
	if err != nil {
		return err
	}

	return nil
}

//check to make sure the required fields of a user are present.
func validateUser(u *user) error {
	if strings.EqualFold(u.Name, "") {
		return errors.New("No name provided.")
	}
	return nil
}
