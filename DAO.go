package main

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Insert the user provided into the user document
func createUser(toInsert *user) error {
	fmt.Printf("Creating user %s\n", toInsert.Name)

	c, ses, err := getUserCollection()
	if err != nil {
		return err
	}
	defer ses.Close()

	err = c.Insert(toInsert)

	if err == nil {
		fmt.Printf("Successfully created user.\n")
	} else {
		fmt.Printf("Error creating the user: %s\n", err.Error())
	}
	return err
}

//Get a user by the hexID string.
func getUserByID(userID string) (user, error) {
	fmt.Printf("Getting User by ID\n")
	var toReturn user
	if !bson.IsObjectIdHex(userID) {
		return toReturn, errors.New("Invalid userID")
	}

	c, ses, err := getUserCollection()
	if err != nil {
		fmt.Printf("Error. %s\n", err.Error())
		return toReturn, err
	}
	defer ses.Close()
	fmt.Printf("Testing: \n%s\n", bson.M{"_id": bson.ObjectIdHex(userID)})
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(userID)}).One(&toReturn)
	if err != nil {
		fmt.Printf("Error. %s\n", err.Error())
		return toReturn, err
	}
	if strings.EqualFold(toReturn.ID, "") {
		err := errors.New("Could not retrieve user: Invalid userID")
		fmt.Printf("Error: %s\n", err.Error())
		return toReturn, err
	}
	fmt.Printf("Done.\n")
	return toReturn, nil
}

func setUserStatus(userID string, status string) error {
	fmt.Printf("Setting user status...\n")
	c, ses, err := getUserCollection()
	if err != nil {
		return err
	}
	defer ses.Close()

	err = c.Update(bson.M{"_id": bson.ObjectIdHex(userID)}, bson.M{
		"$set": bson.M{"status": status}})
	if err == nil {
		fmt.Printf("Done.")
	} else {
		fmt.Printf("Error. %s\n", err.Error())
	}
	return err
}

//Returns a collection, session, and error. Returning the session so that
//close can be called after.
func getUserCollection() (*mgo.Collection, *mgo.Session, error) {
	var sesToReturn *mgo.Session
	var colToReturn *mgo.Collection

	sesToReturn, err := mgo.Dial(config.MongoDBAddress)
	if err != nil {
		return colToReturn, sesToReturn, err
	}

	sesToReturn.SetMode(mgo.Monotonic, true)
	colToReturn = sesToReturn.DB(config.MongoTimeclockDBName).C(config.MongoUserCollectionName)

	return colToReturn, sesToReturn, nil
}
