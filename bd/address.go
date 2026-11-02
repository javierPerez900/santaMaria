package bd

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/javier/santaMaria/models"
)

func InsertAddress(addr models.Address, User string) error {
	fmt.Println("Comienza el Registro InsertAddress")

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "INSERT INTO addresses (Add_UserId, Add_Address, Add_City, Add_State, Add_PostalCode, Add_Phone, Add_Title, Add_Name) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	_, err = Db.Exec(sentencia, User, addr.AddAddress, addr.AddCity, addr.AddState, addr.AddPostalCode, addr.AddPhone, addr.AddTitle, addr.AddName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Insert Address > Ejecución Exitosa")
	return nil
}

func AddressExists(User string, id int) (error, bool) {
	fmt.Println("Comienza AddressExists")

	err := DbConnect()
	if err != nil {
		return err, false
	}
	defer Db.Close()

	var valor int
	err = Db.QueryRow("SELECT 1 FROM addresses WHERE Add_Id = ? AND Add_UserId = ?", id, User).Scan(&valor)
	if err == sql.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return err, false
	}

	fmt.Println("AddressExists > Ejecución Exitosa")
	return nil, true
}

func UpdateAddress(addr models.Address, User string) error {
	fmt.Println("Comienza UpdateAddress")

	partes := []string{}
	args := []any{}

	if addr.AddAddress != "" {
		partes = append(partes, "Add_Address = ?")
		args = append(args, addr.AddAddress)
	}
	if addr.AddCity != "" {
		partes = append(partes, "Add_City = ?")
		args = append(args, addr.AddCity)
	}
	if addr.AddName != "" {
		partes = append(partes, "Add_Name = ?")
		args = append(args, addr.AddName)
	}
	if addr.AddPhone != "" {
		partes = append(partes, "Add_Phone = ?")
		args = append(args, addr.AddPhone)
	}
	if addr.AddPostalCode != "" {
		partes = append(partes, "Add_PostalCode = ?")
		args = append(args, addr.AddPostalCode)
	}
	if addr.AddState != "" {
		partes = append(partes, "Add_State = ?")
		args = append(args, addr.AddState)
	}
	if addr.AddTitle != "" {
		partes = append(partes, "Add_Title = ?")
		args = append(args, addr.AddTitle)
	}

	if len(partes) == 0 {
		return fmt.Errorf("no hay campos para actualizar")
	}

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "UPDATE addresses SET " + strings.Join(partes, ", ") + " WHERE Add_Id = ? AND Add_UserId = ?"
	args = append(args, addr.AddId, User)

	_, err = Db.Exec(sentencia, args...)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Address > Ejecución Exitosa")
	return nil
}

func DeleteAddress(id int, User string) error {
	fmt.Println("Comienza DeleteAddress")

	err := DbConnect()
	if err!=nil {
		return err
	}
	defer Db.Close()

	_, err = Db.Exec("DELETE FROM addresses WHERE Add_Id = ? AND Add_UserId = ?", id, User)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete Address > Ejecución Exitosa")
	return nil
}

func SelectAddress(User string) ([]models.Address, error) {
	fmt.Println("Comienza SelectAddress")

	Addr := []models.Address{}

	err := DbConnect()
	if err!=nil {
		return Addr, err
	}
	defer Db.Close()

	sentencia := "SELECT Add_Id, Add_Address, Add_City, Add_State, Add_PostalCode, Add_Phone, Add_Title, Add_Name FROM addresses WHERE Add_UserId = ?"

	rows, err := Db.Query(sentencia, User)
	if err != nil {
		fmt.Println(err.Error())
		return Addr, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Address
		var addId sql.NullInt64
		var addAddress sql.NullString
		var addCity sql.NullString
		var addState sql.NullString
		var addPostalCode sql.NullString
		var addPhone sql.NullString
		var addTitle sql.NullString
		var addName sql.NullString

		err := rows.Scan(&addId, &addAddress, &addCity, &addState, &addPostalCode, &addPhone, &addTitle, &addName)
		if err != nil {
			return Addr, err
		}

		a.AddId = int(addId.Int64)
		a.AddAddress = addAddress.String
		a.AddCity = addCity.String
		a.AddState = addState.String
		a.AddPostalCode = addPostalCode.String
		a.AddPhone = addPhone.String
		a.AddTitle = addTitle.String
		a.AddName = addName.String

		Addr = append(Addr, a)
	}

	if err := rows.Err(); err != nil {
		return Addr, err
	}

	fmt.Println("Select Addresses > Ejecución Exitosa")
	return Addr, nil
}