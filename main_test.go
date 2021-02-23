package main

import (
	"log"
	"os"
	"testing"

	"fmt"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var a App

func TestMain(m *testing.M) {
    a.Initialize(
        os.Getenv("APP_DB_USERNAME"),
        os.Getenv("APP_DB_PASSWORD"),
        os.Getenv("APP_DB_NAME"))

    ensureTableExists()
    code := m.Run()
    clearTable()
    os.Exit(code)
}

func ensureTableExists() {
    if _, err := a.DB.Exec(tableCreationQuery); err != nil {
        log.Fatal(err)
    }
}

func clearTable() {
    a.DB.Exec("DELETE FROM users")
    a.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
(
    user_id SERIAL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (user_id)
)`


func TestEmptyTable(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/users", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}


func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}


func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}


func TestGetNonExistentUser(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/user/11", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "User not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
    }
}


func TestCreateUser(t *testing.T) {

    clearTable()

    var jsonStr = []byte(`{"username":"testuser", "password": "testpassword"}`)
    req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    response := executeRequest(req)
    checkResponseCode(t, http.StatusCreated, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["username"] != "testuser" {
        t.Errorf("Expected user username to be 'testuser'. Got '%v'", m["username"])
    }

    if m["password"] != "testpassword" {
        t.Errorf("Expected user password to be 'testpassword'. Got '%v'", m["password"])
    }

    // the id is compared to 1.0 because JSON unmarshaling converts numbers to
    // floats, when the target is a map[string]interface{}
    if m["user_id"] != 1.0 {
        t.Errorf("Expected user ID to be '1'. Got '%v'", m["user_id"])
    }
}


func TestGetUser(t *testing.T) {
    clearTable()
    addUsers(3)

    req, _ := http.NewRequest("GET", "/user/1", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)
}

func addUsers(count int) {
    if count < 1 {
        count = 1
    }

    for i := 0; i < count; i++ {
		fmt.Println("INSERT INTO users(username, password) VALUES("+"username_"+strconv.Itoa(i)+", password_"+strconv.Itoa(i)+")")
        a.DB.Exec("INSERT INTO users(username, password) VALUES("+"'username_"+strconv.Itoa(i)+"'"+", 'password_"+strconv.Itoa(i)+"'"+")")
	}
}


func TestUpdateUser(t *testing.T) {

    clearTable()
    addUsers(1)

    req, _ := http.NewRequest("GET", "/user/1", nil)
    response := executeRequest(req)
    var originalUser map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &originalUser)

    var jsonStr = []byte(`{"username":"testUser - updated username", "password": "testpassword"}`)
    req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["user_id"] != originalUser["user_id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["user_id"], m["user_id"])
    }

    if m["username"] == originalUser["username"] {
        t.Errorf("Expected the username to change from '%v' to '%v'. Got '%v'", originalUser["username"], m["username"], m["username"])
    }

    if m["password"] == originalUser["password"] {
        t.Errorf("Expected the password to change from '%v' to '%v'. Got '%v'", originalUser["password"], m["password"], m["password"])
    }
}


func TestDeleteUser(t *testing.T) {
    clearTable()
    addUsers(1)

    req, _ := http.NewRequest("GET", "/user/1", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("DELETE", "/user/1", nil)
    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("GET", "/user/1", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)
}