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

type NullInt sql.NullInt64
type NullString sql.NullString

type User struct {
	ID       NullInt `json:"id"`
	Username NullString `json:"username"`
	Password NullString `json:"password"`
	Email    NullString `json:"email"`
}

func (s *NullInt) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return json.Marshal(s.Int64)
}

func (s *NullInt) Scan(value interface{}) error {
	var v sql.NullInt64
	if err := v.Scan(value); err != nil {
		log.Fatal(err)
		return err
	}
	*s = NullInt(v)
	return nil
}

func (s *NullInt) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.Int64)
	s.Valid = true
	if err != nil {
		s.Valid = false
	}
	return err
}

func (s *NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return json.Marshal("")
	}
	return json.Marshal(s.String)
}

func (s *NullString) Scan(data interface{}) error {
	var str sql.NullString
	if err := str.Scan(data); err != nil {
		log.Fatal(err)
		return err
	}
	*s = NullString(str)
	return nil
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.String)
	s.Valid = true
	if err != nil {
		s.Valid = false
	}
	return err
}

type DB struct {
	db *sql.DB
}

func (u User) String() string {
	return fmt.Sprintln(u.ID, u.Username, u.Password, u.Email)
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
