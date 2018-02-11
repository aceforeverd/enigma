package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type User struct {
	ID       sql.NullInt64
	Username sql.NullString
	Password sql.NullString
	Email    sql.NullString
}

type JsonUser struct {
	ID       int64
	Username string
	Password string
	Email    string
}

func (u JsonUser) User() User {
	return User{
		ID:       sql.NullInt64{Int64: u.ID, Valid: u.ID >= 0},
		Username: sql.NullString{String: u.Username, Valid: len(u.Username) > 0},
		Password: sql.NullString{String: u.Password, Valid: len(u.Password) > 0},
		Email:    sql.NullString{String: u.Email, Valid: len(u.Email) > 0},
	}
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(JsonUser{
		ID:       u.ID.Int64,
		Username: u.Username.String,
		Password: u.Password.String,
		Email:    u.Email.String,
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	var ju JsonUser
	if err := json.Unmarshal(data, &ju); err != nil {
		return err
	}
	*u = ju.User()
	return nil
}

type Users []User

type DB struct {
	db *sql.DB
}

func (u User) String() string {
	return fmt.Sprintln(u.ID, u.Username, u.Password, u.Email)
}

func (u User) JsonEncode() []byte {
	bytes, err := json.Marshal(u)
	if err != nil {
		return []byte{}
	}
	return bytes
}

func (d *DB) init(driveName string, dataSource string) error {
	if d.db == nil {
		db, err := sql.Open(driveName, dataSource)
		if err != nil {
			return err
		}

		d.db = db
		return nil
	}
	return errors.New("db already init")
}

func main() {
	port := flag.String("port", ":8080", "running port")
	flag.Parse()

	var db = &DB{}
	err := db.init("mysql", "test:test@/TEST")
	if err != nil {
		log.Fatal("openning DB", err)
	}

	rows, err := db.db.Query("SELECT * FROM user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	fmt.Println(users)

	router := gin.Default()

	router.GET("/name/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "hello %s", name)
	})

	router.GET("/user", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.Run(*port)
}
