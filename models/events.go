package models

import (
	"errors"
	"example.com/rest-apis/db"
	"fmt"
	"time"
)

type Event struct {
	ID          int64     `json:"id" `
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	DateTime    time.Time `json:"dateTime"`
	UserID      int64     `json:"user_id"`
}

var events = []Event{}

func (e Event) Save() (Event, error) {
	query := `
	INSERT INTO events(name, description, location, dateTime, user_id) 
	VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return e, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return e, err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return e, err
}

func (e Event) Update() error {
	query := `
		UPDATE events
		SET name = ?, description = ?, location = ?, dateTime = ?
		WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	return err
}

func (e Event) Delete() error {
	query := `DELETE  FROM events WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID)
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	var events []Event
	defer rows.Close()
	for rows.Next() {
		var event Event
		readErr := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if readErr != nil {
			return nil, readErr
		}

		events = append(events, event)
	}
	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := "SELECT * FROM events where id = ?"
	row := db.DB.QueryRow(query, id)
	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e Event) Register(userId int64) error {
	query := "INSERT into registations(event_id, user_id) VALUES (?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err
}

func (e Event) CancelRegistation(userId int64) error {
	query := "DELETE FROM registations WHERE event_id = ?  AND user_id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	fmt.Println("Event Id ", e.ID)
	fmt.Println("User Id ", userId)
	result, err := stmt.Exec(e.ID, userId)
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("No event with the id and userID")
	}
	return err
}
