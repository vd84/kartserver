package main

import (
	"database/sql"
	"fmt"
	"strconv"

)

type user struct {
    ID    int     `json:"id"`
    Username  string  `json:"username"`
    Password string `json:"password"`
}

func (p *user) getUser(db *sql.DB) error {
	fmt.Println(db.QueryRow("SELECT username, password FROM users WHERE id="+strconv.Itoa(p.ID)).Scan())
    return db.QueryRow("SELECT username, password FROM users WHERE id="+strconv.Itoa(p.ID)).Scan(&p.Username, &p.Password)
}

func (p *user) updateUser(db *sql.DB) error {
    _, err :=
        db.Exec("UPDATE user SET username=$1, password=$2 WHERE id=$3",
            p.Username, p.Password, p.ID)

    return err
}

func (p *user) deleteUser(db *sql.DB) error {
    _, err := db.Exec("DELETE FROM users WHERE id=$1", p.ID)

    return err
}

func (p *user) createUser(db *sql.DB) error {

	fmt.Println("USER!!!!!!!!!!!!!!! " + p.Username)



    err := db.QueryRow(
        "INSERT INTO users(username, password) VALUES($1, $2) RETURNING id",
        p.Username, p.Password).Scan(&p.ID)

    if err != nil {
        return err
    }

    return nil
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
    rows, err := db.Query(
        "SELECT id, username, password FROM users LIMIT $1 OFFSET $2",
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