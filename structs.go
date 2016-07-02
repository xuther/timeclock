package main

import "time"

type configuration struct {
	DBAddress string
}

type timePunch struct {
	PID         int64
	UID         int64
	In          time.Time
	Out         time.Time
	Words       int64
	Description string
}

type user struct {
	Name    string
	ID      int64
	Status  bool
	punches []timePunch
}
