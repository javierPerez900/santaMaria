package bd

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/javier/santaMaria/models"
	"github.com/javier/santaMaria/secretm"
)

var SecretModel models.SecretRDSJson
var err error
var Db *sql.DB

func ReadSecret() error{
	SecretModel, err = secretm.GetSecret(os.Getenv("SecretName"))
	return err
}

func DbConnect() error {
	Db, err = sql.Open("mysql", ConnStr(SecretModel))
	if err != nil {
		fmt.Println( err.Error())
		return err
	}

	err = Db.Ping()
	if err != nil {
		fmt.Println( err.Error())
		Db.Close()
		return err
	}

	fmt.Println("Conexión exitosa de la BD")
	return nil
}

func ConnStr(claves models.SecretRDSJson) string {
	var dbUser, authToken, dbEndPoint, dbName string
	dbUser = claves.Username
	authToken = claves.Password
	dbEndPoint = claves.Host
	dbName = "santa_maria"
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPasswords=true&parseTime=true", dbUser, authToken, dbEndPoint, dbName)
}

func UserIsAdmin(userUUID string) (bool, string) {
	fmt.Println("Comienza UserIsAdmin")

	err := DbConnect()
	if err != nil {
		return false, err.Error()
	}
	defer Db.Close()

	var valor int
	err = Db.QueryRow("SELECT 1 FROM users WHERE User_UUID = ? AND User_Status = 0", userUUID).Scan(&valor)
	if err == sql.ErrNoRows {
		return false, "User is not admin"
	}
	if err != nil {
		return false, err.Error()
	}

	fmt.Println("UserIsAdmin > Ejecución exitosa")
	return true, ""
}

func UserExists(UserUUID string) (error, bool) {
	fmt.Println("Comienza UserExists")

	err := DbConnect()
	if err != nil {
		return err, false
	}
	defer Db.Close()

	var valor int
	err = Db.QueryRow("SELECT 1 FROM users WHERE User_UUID = ?", UserUUID).Scan(&valor)
	if err == sql.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return err, false
	}

	fmt.Println("UserExists > Ejecución exitosa")
	return nil, true
}