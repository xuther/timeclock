package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func getPunchByID(id int64) (timePunch, error) {
	fmt.Printf("Getting punch %v ...", id)

	punches, err := executePunchQuery("SELECT * FROM punches WHERE pid=?", id)
	if err != nil {
		return timePunch{}, err
	}

	fmt.Printf("Done.")
	//there will be only one, since we're going by ID
	return punches[0], nil
}

func getLastTimepunch(userID int64) (timePunch, error) {
	fmt.Printf("Getting the last punch.\n")

	punches, err := executePunchQuery(
		"SELECT * FROM punches WHERE pid=( SELECT max(pid) FROM punches WHERE uid=?)",
		userID)
	if err != nil {
		return timePunch{}, err
	}

	fmt.Printf("Done\n")
	return punches[0], nil
}

func updatePunch(punch timePunch) error {
	fmt.Printf("Updating Punch...\n")
	if punch.PID == 0 {
		return errors.New("Could not update punch, invalid PID provided.")
	}

	val, err := executeStatement("UPDATE punches SET 'uid'=?, 'in'=?, 'out'=?, 'duration'=?, 'words'=?, 'description'=?"+
		" WHERE pid=? ",
		punch.UID,
		punch.In,
		punch.Out,
		punch.Duration,
		punch.Words,
		punch.Description,
		punch.PID)

	if err != nil {
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

	val, err := executeStatement(
		"INSERT INTO punches('uid', 'in', 'out', 'duration', 'words', 'description') values (?, ?, ?, ?, ?, ?)",
		punch.UID,
		punch.In,
		punch.Out,
		punch.Duration,
		punch.Words,
		punch.Description)

	if err != nil {
		return err
	}
	if val < 1 {
		fmt.Printf("Record not inserted\n.")
		return errors.New("Record not inserted\n.")
	}

	fmt.Printf("Done.\n")
	return nil
}

func executeStatement(statement string, params ...interface{}) (int64, error) {
	db, err := sql.Open("sqlite3", "file:"+config.DBAddress+"?cache=shared&mode=rwc")
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return 0, err
	}
	defer db.Close()

	fmt.Printf("params: %v", params)

	stmt, err := db.Prepare(statement)
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return 0, err
	}

	res, err := stmt.Exec(params...)
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return 0, err
	}

	val, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return 0, err
	}

	return val, nil
}

func executePunchQuery(query string, params ...interface{}) ([]timePunch, error) {
	fmt.Printf("Executing Query %s...", query)

	db, err := sql.Open("sqlite3", "file:"+config.DBAddress+"?cache=shared&mode=rwc")
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Printf("%s\n.", err.Error())
		return nil, err
	}

	res, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}

	punches, err := checkNullPunch(res)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Done.")
	return punches, nil
}

//function to check for the results of a query on the punches table
func checkNullPunch(res *sql.Rows) ([]timePunch, error) {
	var toReturn []timePunch

	for res.Next() {
		var pid sql.NullInt64
		var uid sql.NullInt64
		var in pq.NullTime
		var out pq.NullTime
		var duration sql.NullInt64
		var words sql.NullInt64
		var description sql.NullString

		err := res.Scan(
			&pid,
			&uid,
			&in,
			&out,
			&duration,
			&words,
			&description)

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
		if duration.Valid {
			t.Duration = time.Duration(duration.Int64)
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
