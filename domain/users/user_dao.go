package users

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/oleksandr-ridkodub/bookstore_users-api/datasources/mysql/users_db"
	"github.com/oleksandr-ridkodub/bookstore_users-api/utils/date_utils"
	"github.com/oleksandr-ridkodub/bookstore_users-api/utils/errors"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	errorNoRows      = "no rows in result set"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser     = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
)

func (user *User) Get() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(
				fmt.Sprintf("user %d not found", user.Id))
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error when trying to get user %d: %s", user.Id, err.Error()))
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if saveErr != nil {
		sqlErr, ok := saveErr.(*mysql.MySQLError)
		if !ok {
			return errors.NewInternalServerError(fmt.Sprintf("error while trying to save user: %s", err.Error()))
		}
		fmt.Println(sqlErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		errors.NewInternalServerError(
			fmt.Sprintf("error while trying to save user: %s", err.Error()))
	}
	user.Id = userId

	return nil
}
