package models

import (
	"example.com/rest-apis/db"
	"time"
)

type Event struct {
	ID          int64     `json:"id" `
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	DateTime    time.Time `json:"dateTime"`
	UserID      int       `json:"user_id"`
}

var events = []Event{}

func (e Event) Save() error {
	query := `INSERT INTO events(name, description, location, dateTime, user_id) 
			  VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return nil
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
