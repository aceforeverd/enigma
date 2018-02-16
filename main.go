package main

import (
	"database/sql"
	"flag"
	"github.com/aceforeverd/enigma/controller"
	"github.com/aceforeverd/enigma/repository"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
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

	var ctr controller.UserCon
	ctr = &controller.UserController{Repo: userRepo}

	router := gin.Default()

	router.GET("/users", ctr.GetAll)
	router.GET("/user", ctr.GetUser)
	router.POST("user", ctr.Save)
	router.PUT("/user", ctr.Update)
	router.DELETE("/user", ctr.Delete)

	router.Run(*port)
}
