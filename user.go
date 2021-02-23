package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID       int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (userParam *user) getUser(db *sql.DB) error {
	return db.QueryRow("SELECT username, password FROM users WHERE user_id="+strconv.Itoa(userParam.ID)).Scan(&userParam.Username, &userParam.Password)
}

func (userParam *user) authenticateUser(db *sql.DB) (int, error) {

	fmt.Println("user password gotten from request: " + userParam.Password)
	fmt.Println("user password gotten from request encrypted: " + string(hashPassword(userParam.Password)))

	var user user

	row := db.QueryRow(
		"SELECT user_id, password FROM users where username=$1", userParam.Username)

	err := row.Scan(&user.ID, &user.Password)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return -9999, err
	case nil:
		fmt.Println(user)

	default:
		panic(err)
	}

	fmt.Println("user password stored in db: " + user.Password)

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userParam.Password)); err != nil {
		return -9999, err
	}

	return user.ID, err

}

func (userParam *user) updateUser(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE user SET username=$1, password=$2 WHERE user_id=$3",
			userParam.Username, userParam.Password, userParam.ID)

	return err
}

func (userParam *user) deleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE user_id=$1", userParam.ID)

	return err
}

func (userParam *user) createUser(db *sql.DB) error {

	hashedPassword := hashPassword(userParam.Password)

	err := db.QueryRow(
		"INSERT INTO users(username, password) VALUES($1, $2) RETURNING user_id",
		userParam.Username, hashedPassword).Scan(&userParam.ID)

	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) []byte {

	passwordByteArr := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordByteArr, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))

	err = bcrypt.CompareHashAndPassword(hashedPassword, passwordByteArr)

	return hashedPassword
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	rows, err := db.Query(
		"SELECT user_id, username, password FROM users LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var p user
		if err := rows.Scan(&p.ID, &p.Username, &p.Password); err != nil {
			return nil, err
		}
		users = append(users, p)
	}

	return users, nil
}

func (userParam *user) getUserByName(db *sql.DB) (user, error) {

	row := db.QueryRow(
		"SELECT user_id, username FROM users WHERE username=$1", userParam.Username)

	var dbUser user
	err := row.Scan(&dbUser.ID, &dbUser.Username)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return *userParam, err
	case nil:
		fmt.Println(dbUser)

	default:
		panic(err)
	}
	return dbUser, nil

}
