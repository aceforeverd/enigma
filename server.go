package main

import (
	"database/sql"
	"flag"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", ":8080", "running port")
	flag.Parse()

	db, err := sql.Open("mysql", "test:test@/TEST")
	if err != nil {
		log.Println("openning DB", err)
	}

	if err := db.Ping(); err != nil {
		log.Println("ping: ", err)
	}

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id       sql.NullInt64
			email    sql.NullString
			username sql.NullString
			password sql.NullString
		)
		if err := rows.Scan(&id, &email, &username, &password); err != nil {
			log.Fatal(err)
		}
		log.Println(id, username, password, email)
	}

	router := gin.Default()

	router.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "hello %s", name)
	})

	router.Run(*port)
}
