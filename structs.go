package main

import "time"

type configuration struct {
	MongoDBAddress       string
	MongoTimeclockDBName string
}

type timePunch struct {
	In  time.Time
	Out time.Time
}

type user struct {
	Name    string
	ID      string
	punches []timePunch
}
