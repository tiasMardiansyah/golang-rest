package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type user struct {
	User_id  string `json:"id"`
	Nama     string `json:"nama"`
	Username string `json:"username"`
}

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

	fmt.Println("Connected")

	router := gin.Default()
	router.GET("/user", func(ctx *gin.Context) {
		userList, err := getUser()

		if err != nil {
			return
		}

		ctx.IndentedJSON(http.StatusOK, userList)
	})

	router.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userList, err := getUserById(id)
		if err != nil {
			return
		}

		for _, a := range userList {
			if a.User_id == id {
				ctx.IndentedJSON(http.StatusFound, a)
				return
			}
		}

		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan"})
	})

	router.DELETE("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		//ambil data dari database
		userList, err := getUserById(id)
		if err != nil {
			return
		}

		for i, a := range userList {
			if a.User_id == id {
				userList = append(userList[:i], userList[i+1:]...)
				ctx.IndentedJSON(http.StatusOK, a)
				return
			}
		}

		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan"})
	})

	router.Run(":8080")
}

func getUser() ([]user, error) {
	var users []user

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		return nil, fmt.Errorf("kesalahan mengambil data")
	}

	defer rows.Close()

	for rows.Next() {
		var usr user
		if err := rows.Scan(&usr.User_id, &usr.Nama, &usr.Username); err != nil {
			return nil, fmt.Errorf("GetUser: %v", err)
		}

		users = append(users, usr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUser: %v", err)
	}

	return users, nil
}

func getUserById(id string) ([]user, error) {
	var users []user

	rows, err := db.Query("SELECT * FROM user where id=?", id)
	if err != nil {
		return nil, fmt.Errorf("Kesalahan mengambil data")
	}
	defer rows.Close()

	for rows.Next() {
		var usr user
		if err := rows.Scan(&usr.User_id, &usr.Nama, &usr.Username); err != nil {
			return nil, fmt.Errorf("GetUser: %v", err)
		}

		users = append(users, usr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUser: %v", err)
	}

	return users, nil
}
