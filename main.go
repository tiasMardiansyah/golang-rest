package main

import (
	"database/sql"
	"example/web-service-gin/database/queries"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {

	//isinya configurasi sql
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "latihan_golang",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to Database: ", cfg.DBName)

	router := gin.Default()
	router.GET("/user", func(ctx *gin.Context) {
		userList, statusCode, err := queries.GetUser()
		if err != nil {
			ctx.IndentedJSON(statusCode, gin.H{"message": http.StatusText(statusCode)})
			return
		}

		ctx.IndentedJSON(statusCode, userList)
	})

	router.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userList, statusCode, err := queries.GetUserById(id)
		if err != nil {
			ctx.IndentedJSON(statusCode, gin.H{"message": http.StatusText(statusCode)})
			return
		}

		ctx.IndentedJSON(statusCode, userList)
	})

	router.DELETE("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		statusCode, err := queries.DeleteUser(id)
		if err != nil {
			return
		}

		ctx.IndentedJSON(statusCode, gin.H{"message": http.StatusText(statusCode)})
	})

	router.Run(":8080")
}
