package bd

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/javier/santaMaria/models"
	"github.com/javier/santaMaria/tools"
)

func UpdateUser(UField models.User, User string) error {
	fmt.Println("Comienza UpdateUser")

	partes := []string{}
	args := []any{}

	// Los nombres de columna son texto fijo; solo los valores van como parámetros
	if len(UField.UserFirstName) > 0 {
		partes = append(partes, "User_FirstName = ?")
		args = append(args, UField.UserFirstName)
	}
	if len(UField.UserLastName) > 0 {
		partes = append(partes, "User_LastName = ?")
		args = append(args, UField.UserLastName)
	}

	if len(partes) == 0 {
		return fmt.Errorf("no hay campos para actualizar")
	}

	partes = append(partes, "User_DateUpg = ?")
	args = append(args, tools.FechaMySQL())

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "UPDATE users SET " + strings.Join(partes, ", ") + " WHERE User_UUID = ?"
	args = append(args, User)

	_, err = Db.Exec(sentencia, args...)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update User > Ejecución Exitosa")
	return nil
}

func SelectUser(UserId string) (models.User, error) {
	fmt.Println("Comienza SelectUser")
	User := models.User{}

	err := DbConnect()
	if err != nil {
		return User, err
	}
	defer Db.Close()

	var firstName sql.NullString
	var lastName sql.NullString
	var dateUpg sql.NullTime

	err = Db.QueryRow("SELECT * FROM users WHERE User_UUID = ?", UserId).Scan(
		&User.UserUUID, &User.UserEmail, &firstName, &lastName, &User.UserStatus, &User.UserDateAdd, &dateUpg,
	)
	if err != nil {
		fmt.Println(err.Error())
		return User, err
	}

	User.UserFirstName = firstName.String
	User.UserLastName = lastName.String
	User.UserDateUpd = dateUpg.Time.String()

	fmt.Println("Select User > Ejecución Exitosa")
	return User, nil
}

func SelectUsers(Page int) (models.ListUsers, error) {
	fmt.Println("Comienza SelectUsers")

	var lu models.ListUsers
	User := []models.User{}

	err := DbConnect()
	if err != nil {
		return lu, err
	}
	defer Db.Close()

	if Page < 1 {
		Page = 1
	}
	offset := (Page - 1) * 10

	var registros int
	err = Db.QueryRow("SELECT count(*) FROM users").Scan(&registros)
	if err != nil {
		return lu, err
	}
	lu.TotalItems = registros

	rows, err := Db.Query("SELECT * FROM users LIMIT ? OFFSET ?", 10, offset)
	if err != nil {
		fmt.Println(err.Error())
		return lu, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		var firstName sql.NullString
		var lastName sql.NullString
		var dateUpg sql.NullTime

		err := rows.Scan(&u.UserUUID, &u.UserEmail, &firstName, &lastName, &u.UserStatus, &u.UserDateAdd, &dateUpg)
		if err != nil {
			return lu, err
		}

		u.UserFirstName = firstName.String
		u.UserLastName = lastName.String
		u.UserDateUpd = dateUpg.Time.String()
		User = append(User, u)
	}

	if err := rows.Err(); err != nil {
		return lu, err
	}

	fmt.Println("Select Users > Ejecución Exitosa")

	lu.Data = User
	return lu, nil
}