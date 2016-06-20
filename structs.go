package main

import "time"

type configuration struct {
	MongoDBAddress          string
	MongoTimeclockDBName    string
	MongoUserCollectionName string
}

type timePunch struct {
	In  time.Time
	Out time.Time
}

type user struct {
	Name    string
	ID      string `bson:"_id,omitempty"`
	Status  string
	punches []timePunch
}
