package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

func createUser(toInsert *user) error {
	fmt.Printf("Creating user %s\n", toInsert.Name)

	session, err := mgo.Dial(config.MongoDBAddress)
	if err != nil {
		return err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.MongoTimeclockDBName).C("Users")

	err = c.Insert(toInsert)

	if err == nil {
		fmt.Printf("Successfully created user.\n")
	} else {
		fmt.Printf("Error creating the user: %s\n", err.Error())
	}
	return err
}
