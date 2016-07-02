package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
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

func updatePunch(punch timePunch) error {
	fmt.Printf("Updating Punch...\n")
	if punch.PID == 0 {
		return errors.New("Could not update punch, invalid PID provided.")
	}

	db, err := sql.Open("sqlite3", "file:"+config.DBAddress+"?cache=shared&mode=rwc")
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"UPDATE punches SET 'uid'=?, 'in'=?, 'out'=?, 'words'=?, 'description'=?" +
			" WHERE pid=? ")
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return err
	}

	res, err := stmt.Exec(
		punch.UID,
		punch.In,
		punch.Out,
		punch.Words,
		punch.Description,
		punch.PID)

	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return err
	}

	val, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return err
	}

	if val < 1 {
		fmt.Printf("Record not inserted\n.")
		return errors.New("Record not inserted\n.")
	}
	fmt.Printf("Done.\n")

	return nil
}

func createPunch(punch timePunch) error {
	fmt.Printf("Creating Punch...\n")

	db, err := sql.Open("sqlite3", "file:"+config.DBAddress+"?cache=shared&mode=rwc")
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"INSERT INTO punches('uid', 'in', 'out', 'words', 'description') values (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return err
	}

	res, err := stmt.Exec(punch.UID, punch.In, punch.Out, punch.Words, punch.Description)
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return err
	}

	val, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return err
	}

	if val < 1 {
		fmt.Printf("Record not inserted\n.")
		return errors.New("Record not inserted\n.")
	}
	fmt.Printf("Done.\n")

	return nil
}

//function to check for the results of a query on the punches table
func checkNullPunch(res *sql.Rows) ([]timePunch, error) {
	var toReturn []timePunch

	for res.Next() {
		var pid sql.NullInt64
		var uid sql.NullInt64
		var in pq.NullTime
		var out pq.NullTime
		var words sql.NullInt64
		var description sql.NullString

		err := res.Scan(
			&pid,
			&uid,
			&in,
			&out,
			&words,
			&description)
		fmt.Printf("We're here! %v %v %v\n", pid, uid, in)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}
		t := timePunch{}

		if pid.Valid {
			t.PID = pid.Int64
		}
		if uid.Valid {
			t.UID = uid.Int64
		}
		if words.Valid {
			t.Words = words.Int64
		}
		if description.Valid {
			t.Description = description.String
		}
		if in.Valid {
			t.In = in.Time
		}
		if out.Valid {
			t.Out = out.Time
		}
		fmt.Printf("%+v\n", t)

		toReturn = append(toReturn, t)
	}
	if len(toReturn) < 1 {
		fmt.Printf("No punch found\n")
		return toReturn, errors.New("No punch found\n")
	}
	return toReturn, nil
}

func getLastTimepunch(userID int64) (timePunch, error) {
	fmt.Printf("Getting the last punch.\n")
	db, err := sql.Open("sqlite3", "file:"+config.DBAddress+"?cache=shared&mode=rwc")
	if err != nil {
		return timePunch{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(
		"SELECT * FROM punches WHERE pid=( SELECT max(pid) FROM punches WHERE uid=?)")
	if err != nil {
		return timePunch{}, err
	}

	res, err := stmt.Query(userID)
	if err != nil {
		return timePunch{}, err
	}

	vals, err := checkNullPunch(res)
	if err != nil {
		return timePunch{}, err
	}

	toReturn := vals[0]

	fmt.Printf("Done\n")
	return toReturn, nil
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
