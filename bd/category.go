package bd

import (
	"database/sql"
	"fmt"

	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/javier/santaMaria/models"
)

func InsertCategory(c models.Category) (int64, error) {
	fmt.Println("Comienza Registro de InsertCategory")

	err := DbConnect()
	if err!=nil {
		return 0, err
	}
	defer Db.Close()

	result, err := Db.Exec("INSERT INTO category (Categ_Name, Categ_Path) VALUES (?, ?)", c.CategName, c.CategPath)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	LastInsertId, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	fmt.Println("Insert Category > Ejecución Exitosa")
	return LastInsertId, nil
}

func UpdateCategory(c models.Category) error {
	fmt.Println("Comienza Registro de UpdateCategory")

	partes := []string{}
	args := []any{}

	if len(c.CategName) > 0 {
		partes = append(partes, "Categ_Name = ?")
		args = append(args, c.CategName)
	}
	if len(c.CategPath) > 0 {
		partes = append(partes, "Categ_Path = ?")
		args = append(args, c.CategPath)
	}

	if len(partes) == 0 {
		return fmt.Errorf("no hay campos para actualizar")
	}

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "UPDATE category SET " + strings.Join(partes, ", ") + " WHERE Categ_Id = ?"
	args = append(args, c.CategID)

	_, err = Db.Exec(sentencia, args...)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Category > Ejecución Exitosa")
	return nil
}

func DeleteCategory(id int) error {
	fmt.Println("Comienza Registro de DeleteCategory")

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	_, err = Db.Exec("DELETE FROM category WHERE Categ_Id = ?", id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete Category > Ejecución Exitosa")
	return nil
}

func SelectCategories(CategId int, Slug string) ([]models.Category, error) {
	fmt.Println("Comienza SelectCategories")

	Categ := []models.Category{}

	err := DbConnect()
	if err != nil {
		return Categ, err
	}
	defer Db.Close()

	sentencia := "SELECT Categ_Id, Categ_Name, Categ_Path FROM category"
	args := []any{}

	if CategId > 0 {
		sentencia += " WHERE Categ_Id = ?"
		args = append(args, CategId)
	} else if len(Slug) > 0 {
		sentencia += " WHERE Categ_Path LIKE ?"
		args = append(args, "%"+Slug+"%")
	}

	rows, err := Db.Query(sentencia, args...)
	if err != nil {
		return Categ, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Category
		var categId sql.NullInt32
		var categName sql.NullString
		var categPath sql.NullString

		err := rows.Scan(&categId, &categName, &categPath)
		if err != nil {
			return Categ, err
		}

		c.CategID = int(categId.Int32)
		c.CategName = categName.String
		c.CategPath = categPath.String

		Categ = append(Categ, c)
	}

	if err := rows.Err(); err != nil {
		return Categ, err
	}

	fmt.Println("Select Category > Ejecución Exitosa")
	return Categ, nil
}