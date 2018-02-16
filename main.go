package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aceforeverd/enigma/repository"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

// InitDB initial a database connection
func InitDB(driveName string, dataSource string) (*sql.DB, error) {
	db, err := sql.Open(driveName, dataSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	port := flag.String("port", ":8080", "running port")
	flag.Parse()

	db, err := InitDB("mysql", "test:test@/TEST")
	if err != nil {
		log.Fatal("openning DB", err)
	}
	var userRepo repository.UserRepo
	userRepo = &repository.UserRepoIml{DB: db}
	if err := userRepo.InitTable(); err != nil {
		log.Fatal(err)
	}

	userList, err := userRepo.GetAll()

	fmt.Println("marshal:")
	data, err := json.Marshal(userList)
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
