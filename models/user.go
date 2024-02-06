package models

import (
	"example.com/rest-apis/db"
	"example.com/rest-apis/utils"
	"fmt"
	"github.com/pkg/errors"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?,?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	hashPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(u.Email, hashPassword)
	if err != nil {
		return err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = userID
	return nil
}

func (u User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)
	var retrivedPassword string
	err := row.Scan(&u.ID, &retrivedPassword)
	if err != nil {
		return errors.New("Invalid Password")
	}
	if !utils.CheckPasswordHash(u.Password, retrivedPassword) {
		return errors.New("Invalid Password")
	}
	fmt.Println("In validate cred ", u)
	return nil
}
