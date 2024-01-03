package models

import (
	"crypto/sha256"
	"database/sql"
	"time"
)

type User struct {
    Uid int64
    Username string
    Business_name string
    Password string
    Created_at time.Time
    Updated_at time.Time
}

func (user *User) Save() (int64, error) {
    db, err := sql.Open("sqlite3", "database/httpserver.db")
    CheckErr(err)
    defer db.Close()

    var stmt *sql.Stmt

    if (user.Uid > 0) {
        stmt, err = db.Prepare("UPDATE user SET username=(?), business_name=(?), password=(?), updated_at=Date('now')")
        CheckErr(err)
    } else {
        stmt, err = db.Prepare("INSERT INTO user(username, business_name, password, created_at, updated_at) values(?,?,?,DATE('now'), DATE('now'))")
        CheckErr(err)
    }

    res, err := stmt.Exec(user.Username, user.Business_name, hashPassword(user.Password))
    CheckErr(err)

    id, err := res.LastInsertId()
    CheckErr(err)

    return id, nil
}

func LoadUser(id int64) (*User, error) {
    db, err := sql.Open("sqlite3", "database/httpserver.db")
    CheckErr(err)

    defer db.Close()

    rows, err := db.Query("SELECT * FROM user WHERE uid=(?) LIMIT 1", id)
    CheckErr(err)
    defer rows.Close()

    var user User

    for rows.Next() {
        err = rows.Scan(&user.Uid, &user.Username, &user.Password, &user.Business_name, &user.Created_at, &user.Updated_at)
        CheckErr(err)

    }
    return &user, nil
}

func LoginUser(username string, password string) (*User, error) {
    db, err := sql.Open("sqlite3", "database/httpserver.db")
    CheckErr(err)

    defer db.Close()

    rows, err := db.Query("SELECT * FROM user WHERE username=(?) AND password=(?) LIMIT 1", username, hashPassword(password))
    CheckErr(err)
    defer rows.Close()

    var user User

    for rows.Next() {
        err = rows.Scan(&user.Uid, &user.Username, &user.Password, &user.Business_name, &user.Created_at, &user.Updated_at)
        CheckErr(err)

    }
    return &user, nil
}

func FindUserByPassword(password string) (*User, error) {
    db, err := sql.Open("sqlite3", "database/httpserver.db")
    CheckErr(err)

    defer db.Close()

    rows, err := db.Query("SELECT * FROM user WHERE password=(?) LIMIT 1", hashPassword(password))
    CheckErr(err)
    defer rows.Close()

    var user User

    for rows.Next() {
        err = rows.Scan(&user.Uid, &user.Username, &user.Password, &user.Business_name, &user.Created_at, &user.Updated_at)
        CheckErr(err)
    }
    return &user, nil
}

func hashPassword(password string) (string) {
    hash := sha256.New()
    hash.Write([]byte(password))

    return string(hash.Sum(nil))
}

func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}

func GetUsers(page int) ([]User, error) {
    db, err := sql.Open("sqlite3", "database/httpserver.db")
    CheckErr(err)

    defer db.Close()

    rows, err := db.Query("SELECT * FROM user LIMIT 100 OFFSET (?)", page)
    CheckErr(err)
    defer rows.Close()

    var users []User

    for rows.Next() {
        var user User
        err = rows.Scan(&user.Uid, &user.Username, &user.Password, &user.Business_name, &user.Created_at, &user.Updated_at)
        CheckErr(err)

        users = append(users, user)
    }
    return users, nil
}

//TODO: create a proper login system and proper password storage, as plain text is cringe
