package queries

import (
	"database/sql"
	"example/web-service-gin/database/model"
	"fmt"
	"net/http"
)

const deleteUser = "DELETE FROM user WHERE id=?"
const updateUser = "UPDATE user SET nama=?, username=? WHERE id=?"
const createUser = "INSERT INTO user (id,nama,username) VALUES(?,?,?)"

const getAllUsers = "SELECT * FROM user"
const getUserById = "SELECT * FROM user WHERE id=?"

func GetUser(db *sql.DB) ([]model.User, int, error) {
	var users []model.User

	rows, err := db.Query(getAllUsers)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("")
	}

	defer rows.Close()

	for rows.Next() {
		var usr model.User
		if err := rows.Scan(&usr.User_id, &usr.Nama, &usr.Username); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("GetUser(): %v", err)
		}

		users = append(users, usr)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("GetUser: %v", err)
	}

	return users, http.StatusOK, nil
}

func GetUserById(db *sql.DB, id string) ([]model.User, int, error) {
	var users []model.User

	rows, err := db.Query(getUserById, id)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("GetUserById(%s): %v", id, err)
	}
	defer rows.Close()

	for rows.Next() {
		var usr model.User
		if err := rows.Scan(&usr.User_id, &usr.Nama, &usr.Username); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("GetUserById(%s): %v", id, err)
		}
		users = append(users, usr)
	}

	if len(users) == 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("GetUser(%s): %v", id, err)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("GetUser(%s): %v", id, err)
	}

	return users, http.StatusOK, nil
}

func CreateUser(db *sql.DB, user model.User) (int, error) {
	result, err := db.Exec(createUser, nil, user.Nama, user.Username)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("CreateUser(): %v", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("CreateUser(): %v", err)
	}

	if rowAffected == 0 {
		return http.StatusBadRequest, fmt.Errorf("CreateUser(): Bad Request")
	}

	return http.StatusOK, nil
}

func UpdateUser(db *sql.DB, user model.User) (int, error) {
	result, err := db.Exec(updateUser, user.Nama, user.Username, user.User_id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("CreateUser(): %v", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("CreateUser(): %v", err)
	}

	if rowAffected == 0 {
		return http.StatusBadRequest, fmt.Errorf("CreateUser(): Bad Request")
	}

	return http.StatusOK, nil
}

func DeleteUser(db *sql.DB, id string) (int, error) {
	result, err := db.Exec(deleteUser, id)

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("deleteUser(%s): %v", id, err)
	}

	rowAffected, err := result.RowsAffected()
	if rowAffected == 0 {
		return http.StatusBadRequest, fmt.Errorf("deleteUser(%s): %s", id, "Akun tidak ditemukan")
	}

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("deleteUser(%s): %v", id, err)
	}

	return http.StatusOK, nil
}
