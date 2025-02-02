package service

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wishwaprabodha/go-server/connection"
	"github.com/wishwaprabodha/go-server/model"
	"log"
)

type User struct {
	UserId       int    `json:"userId"`
	UserName     string `json:"userName"`
	UserEmail    string `json:"userEmail"`
	UserPassword string `json:"userPassword"`
	ModifiedTime string `json:"modified_time,omitempty"`
}

func CreateUserDeta(u *model.User) (key string, err error) {
	conn, err := connection.DetaConnection()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	key, err2 := conn.Insert(u)
	if err != nil {
		log.Fatal(err2)
		return "", err2
	}
	return key, nil
}

func CreateUser(u *model.User) (error, sql.Result) {
	//<package>.<imported_function_name>
	stmt, err := connection.DBConnection().Prepare("INSERT INTO User(userId, userName, userEmail, userPassword) values (?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Error: ", err)
		return err, nil
	}
	createdRecord, insertError := stmt.Exec(u.UserId, u.UserName, u.UserEmail, u.UserPassword)
	if insertError != nil {
		log.Fatal("Error: ", insertError)
		return insertError, nil
	}
	log.Println(createdRecord)
	return nil, createdRecord
}

func GetUserByEmail(email string) (error, User) {
	queryString := fmt.Sprintf(`SELECT * FROM User WHERE userEmail = '%s'`, email)
	var user User
	err := connection.DBConnection().QueryRow(queryString).Scan(&user.UserId,
		&user.UserName, &user.UserEmail, &user.UserPassword, &user.ModifiedTime)
	fmt.Println(user)
	if err != nil {
		return err, user
	}
	return nil, user
}
