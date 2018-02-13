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
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u User) JSONUser() JsonUser {
	id, _ := u.ID.Value()
	username, _ := u.Username.Value()
	passwd, _ := u.Password.Value()
	var password string
	var ok bool
	if password, ok = passwd.(string); !ok {
		password = ""
	}
	email, _ := u.Email.Value()

	return JsonUser{
		ID:       id.(int64),
		Username: username.(string),
		Password: password,
		Email:    email.(string),
	}
}

func (u JsonUser) User() User {
	var user User
	if err := user.ID.Scan(u.ID); err != nil {
		panic(err)
	}
	if err := user.Username.Scan(u.Username); err != nil {
		panic(err)
	}
	if err := user.Password.Scan(u.Password); err != nil {
		panic(err)
	}
	if err := user.Email.Scan(u.Email); err != nil {
		panic(err)
	}
	return user
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.JSONUser())
}

func (u *User) UnmarshalJSON(data []byte) error {
	var ju JsonUser
	if err := json.Unmarshal(data, &ju); err != nil {
		return err
	}
	*u = ju.User()
	return nil
}

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

	fmt.Println("marshal:")
	data, err := json.Marshal(users)
	if err == nil {
		fmt.Println(string(data))
	} else {
		panic(err)
	}

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
