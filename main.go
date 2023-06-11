package main

import (
	"database/sql"
	"example/web-service-gin/database/model"
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
		userList, statusCode, err := queries.GetUser(db)
		if err != nil {
			ctx.IndentedJSON(statusCode, gin.H{"message": http.StatusText(statusCode)})
			return
		}

		ctx.IndentedJSON(statusCode, userList)
	})

	router.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userList, statusCode, err := queries.GetUserById(db, id)
		if err != nil {
			ctx.IndentedJSON(statusCode, gin.H{"message": http.StatusText(statusCode)})
			return
		}

		ctx.IndentedJSON(statusCode, userList)
	})

	router.POST("/user", func(ctx *gin.Context) {
		nama := ctx.PostForm("nama")
		username := ctx.PostForm("username")

		if nama == "" || username == "" {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
			return
		}

		var newUser = model.User{
			Nama:     nama,
			Username: username,
		}

		statusCode, _ := queries.CreateUser(db, newUser)
		ctx.IndentedJSON(statusCode, gin.H{"message": http.StatusText(statusCode)})
	})

	router.PUT("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		nama := ctx.PostForm("nama")
		username := ctx.PostForm("username")

		if nama == "" || username == "" {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
			return
		}

		var updatedUser = model.User{
			User_id:  id,
			Nama:     nama,
			Username: username,
		}

		statusCode, _ := queries.UpdateUser(db, updatedUser)
		ctx.IndentedJSON(statusCode, gin.H{"message": http.StatusText(statusCode)})

	})

	router.DELETE("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		statusCode, err := queries.DeleteUser(db, id)
		if err != nil {
			return
		}

		ctx.IndentedJSON(statusCode, gin.H{"message": http.StatusText(statusCode)})
	})

	router.Run(":8080")
}
