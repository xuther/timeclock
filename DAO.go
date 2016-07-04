package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

//Insert the user provided into the user document
func createUser(toInsert *user) error {
	fmt.Printf("Creating User...\n")
	db, err := sql.Open("sqlite3", config.DBAddress)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"INSERT INTO users (username, password, status) values(?,?,?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(toInsert.Name, "", toInsert.Status)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	toInsert.ID = id

	fmt.Printf("Done.")
	return nil
}

func listtables() error {
	fmt.Printf("Printing tables...")
	db, err := sql.Open("sqlite3", config.DBAddress)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"SELECT tbl_name FROM sqlite_master WHERE type = ?")
	if err != nil {
		return err
	}

	res, err := stmt.Query("table")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", res)
	defer res.Close()

	strings, err := res.Columns()
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", strings)
	var test string
	if res.Next() {
		err = res.Scan(&test)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", test)
	}
	return nil
}

//Get a user by the ID.
func getUserByID(userID int64) (user, error) {
	fmt.Printf("Getting User by ID\n")

	db, err := sql.Open("sqlite3", "file:"+config.DBAddress+"?cache=shared&mode=rwc")
	if err != nil {
		return user{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"SELECT * FROM users WHERE uid = ?")
	if err != nil {
		fmt.Printf("Couldn't prepare query.\n")
		return user{}, err
	}

	res, err := stmt.Query(userID)
	if err != nil {
		fmt.Printf("Couldn't execute query.\n")
		return user{}, err
	}
	defer res.Close()

	for res.Next() {
		var uid int64
		var username string
		var password string
		var status bool

		err = res.Scan(&uid, &username, &password, &status)
		if err != nil {
			return user{}, err
		}

		return user{Name: username, ID: uid, Status: status}, nil
	}

	return user{}, nil
}

func setUserStatus(userID int64, status bool) error {
	fmt.Printf("Setting user status...\n")

	db, err := sql.Open("sqlite3", "file:"+config.DBAddress+"?cache=shared&mode=rwc")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"UPDATE users SET status=? WHERE uid=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(status, userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected < 1 {
		fmt.Printf("No user found %v...\n", userID)
		return errors.New("No user found.\n")
	}

	fmt.Printf("Done.\n")
	return nil
}
