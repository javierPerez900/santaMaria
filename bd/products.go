package bd

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/javier/santaMaria/models"
)

// ErrNoFields se devuelve cuando UpdateProduct no recibe ningún campo para actualizar.
var ErrNoFields = errors.New("no hay campos para actualizar")

// func InsertProduct(p models.Product) (int64,error) {
// 	fmt.Println("Comienza Registro")

// 	err := DbConnect()
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer Db.Close()

// 	sentencia := "INSERT INTO products (Prod_Title "

// 	if len(p.ProdDescription)>0 {
// 		sentencia += ", Prod_Description"
// 	}
// 	if p.ProdPrice > 0 {
// 		sentencia += ", Prod_Price"
// 	}
// 	if p.ProdCategId > 0 {
// 		sentencia += ", Prod_CategoryId"
// 	}
// 	if p.ProdStock > 0 {
// 		sentencia += ", Prod_Stock"
// 	}
// 	if len(p.ProdPath) > 0 {
// 		sentencia += ", Prod_Path"
// 	}
// 	if p.ProdSellOnlyUnit == 0 || p.ProdSellOnlyUnit == 1 {
// 		sentencia += ", Prod_Sell_Only_Unit"

// 		if p.ProdSellOnlyUnit == 0 {
// 			if p.ProdPriceBox > 0 {
// 				sentencia += ", Prod_Price_Box"
// 			}
// 			if p.ProdUnitsPerBox > 0 {
// 				sentencia += ", Prod_Units_Per_Box"
// 			}
// 			if p.ProdWholesaleMinQty > 0 {
// 				sentencia += ", Prod_Wholesale_Min_Qty"
// 			}
// 			if p.ProdPriceBox > 0 && p.ProdUnitsPerBox > 0{
// 				sentencia += ", Prod_Price_Wholesale_Unit"
// 			}
// 		}

// 	}

// 	sentencia += ") VALUES ('" + tools.EscapeString(p.ProdTitle) + "'"

// 	if len(p.ProdDescription)>0 {
// 		sentencia += ", '" + tools.EscapeString(p.ProdDescription) + "'"
// 	}
// 	if p.ProdPrice > 0 {	
// 		sentencia += ", " + strconv.FormatFloat(p.ProdPrice, 'e', -1, 64)
// 	}
// 	if p.ProdCategId > 0 {
// 		sentencia += ", " + strconv.Itoa(p.ProdCategId)
// 	}
// 	if p.ProdStock > 0 {
// 		sentencia += ", " + strconv.Itoa(p.ProdStock)
// 	}
// 	if len(p.ProdPath) > 0 {
// 		sentencia += ", '" + tools.EscapeString(p.ProdPath) + "'"
// 	}
// 	if p.ProdSellOnlyUnit == 0 || p.ProdSellOnlyUnit == 1 {
// 		sentencia += ", " + strconv.Itoa(p.ProdSellOnlyUnit)

// 		if p.ProdSellOnlyUnit == 0 {
// 			if p.ProdPriceBox > 0 {
// 				sentencia += ", " + strconv.FormatFloat(p.ProdPriceBox, 'e', -1, 64)
// 			}
// 			if p.ProdUnitsPerBox > 0 {
// 				sentencia += ", " + strconv.Itoa(p.ProdUnitsPerBox)
// 			}
// 			if p.ProdWholesaleMinQty > 0 {
// 				sentencia += ", " + strconv.Itoa(p.ProdWholesaleMinQty)
// 			}
// 			if p.ProdPriceBox > 0 && p.ProdUnitsPerBox > 0{
// 			priceWholesaleUnit := p.ProdPriceBox / float64(p.ProdUnitsPerBox)
// 				sentencia += ", " + strconv.FormatFloat(priceWholesaleUnit, 'e', -1, 64)
// 			}
// 		}

// 	}

// 	sentencia += ")"

// 	var result sql.Result
// 	result, err = Db.Exec(sentencia)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return 0, err
// 	}

// 	LastInsertId, err2 := result.LastInsertId()
// 	if err2 != nil {
// 		return 0, err2
// 	}

// 	fmt.Println("Insert Product > Ejecución Exitosa")
// 	return LastInsertId, nil
// }

func InsertProduct(p models.Product) (int64, error) {
	fmt.Println("Comienza Registro")

	err := DbConnect()
	if err != nil {
		return 0, err
	}
	defer Db.Close()

	// Los nombres de columna son texto fijo; solo los valores van como parámetros
	columnas := []string{"Prod_Title"}
	args := []any{p.ProdTitle}

	agregar := func(columna string, valor any) {
		columnas = append(columnas, columna)
		args = append(args, valor)
	}

	if len(p.ProdDescription) > 0 {
		agregar("Prod_Description", p.ProdDescription)
	}
	if p.ProdPrice > 0 {
		agregar("Prod_Price", p.ProdPrice)
	}
	if p.ProdCategId > 0 {
		agregar("Prod_CategoryId", p.ProdCategId)
	}
	if p.ProdStock > 0 {
		agregar("Prod_Stock", p.ProdStock)
	}
	if len(p.ProdPath) > 0 {
		agregar("Prod_Path", p.ProdPath)
	}
	if p.ProdSellOnlyUnit == 0 || p.ProdSellOnlyUnit == 1 {
		agregar("Prod_Sell_Only_Unit", p.ProdSellOnlyUnit)

		if p.ProdSellOnlyUnit == 0 {
			if p.ProdPriceBox > 0 {
				agregar("Prod_Price_Box", p.ProdPriceBox)
			}
			if p.ProdUnitsPerBox > 0 {
				agregar("Prod_Units_Per_Box", p.ProdUnitsPerBox)
			}
			if p.ProdWholesaleMinQty > 0 {
				agregar("Prod_Wholesale_Min_Qty", p.ProdWholesaleMinQty)
			}
			if p.ProdPriceBox > 0 && p.ProdUnitsPerBox > 0 {
				agregar("Prod_Price_Wholesale_Unit", p.ProdPriceBox/float64(p.ProdUnitsPerBox))
			}
		}
	}

	marcadores := strings.TrimSuffix(strings.Repeat("?, ", len(columnas)), ", ")
	sentencia := "INSERT INTO products (" + strings.Join(columnas, ", ") + ") VALUES (" + marcadores + ")"

	result, err := Db.Exec(sentencia, args...)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	LastInsertId, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	fmt.Println("Insert Product > Ejecución Exitosa")
	return LastInsertId, nil
}

func UpdateProduct(p models.Product) error {
	fmt.Println("Comienza UpdateProduct")

	partes := []string{}
	args := []any{}

	agregar := func(columna string, valor any) {
		partes = append(partes, columna+" = ?")
		args = append(args, valor)
	}

	// Solo se actualizan los campos enviados (no vacíos ni en 0)
	if p.ProdTitle != "" {
		agregar("Prod_Title", p.ProdTitle)
	}
	if p.ProdDescription != "" {
		agregar("Prod_Description", p.ProdDescription)
	}
	if p.ProdPrice > 0 {
		agregar("Prod_Price", p.ProdPrice)
	}
	if p.ProdCategId > 0 {
		agregar("Prod_CategoryId", p.ProdCategId)
	}
	if p.ProdStock > 0 {
		agregar("Prod_Stock", p.ProdStock)
	}
	if p.ProdPath != "" {
		agregar("Prod_Path", p.ProdPath)
	}
	if p.ProdSellOnlyUnit > 0 {
		agregar("Prod_Sell_Only_Unit", p.ProdSellOnlyUnit)
	}
	if p.ProdPriceBox > 0 {
		agregar("Prod_Price_Box", p.ProdPriceBox)
	}
	if p.ProdUnitsPerBox > 0 {
		agregar("Prod_Units_Per_Box", p.ProdUnitsPerBox)
	}
	if p.ProdWholesaleMinQty > 0 {
		agregar("Prod_Wholesale_Min_Qty", p.ProdWholesaleMinQty)
	}
	if p.ProdPriceBox > 0 && p.ProdUnitsPerBox > 0 {
		agregar("Prod_Price_Wholesale_Unit", p.ProdPriceBox/float64(p.ProdUnitsPerBox))
	}

	if len(partes) == 0 {
		return ErrNoFields
	}

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "UPDATE products SET " + strings.Join(partes, ", ") + " WHERE Prod_Id = ?"
	args = append(args, p.ProdId)

	_, err = Db.Exec(sentencia, args...)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Product > Ejecución Exitosa")
	return nil
}

func DeleteProduct(id int) error {
	fmt.Println("Comienza Registro de DeleteProduct")

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	_, err = Db.Exec("DELETE FROM products WHERE Prod_Id = ?", id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete Product > Ejecución Exitosa")
	return nil
}

func SelectProduct(p models.Product, choice string, page int, pageSize int, orderType string, orderField string) (models.ProductResp, error) {
	fmt.Println("Comienza SelectProduct")
	var Resp models.ProductResp
	Prod := []models.Product{}

	err := DbConnect()
	if err != nil {
		return Resp, err
	}
	defer Db.Close()

	sentencia := "SELECT Prod_Id, Prod_Title, Prod_Description, Prod_CreatedAt, Prod_Updated, Prod_Price, Prod_Price_Box, Prod_Units_Per_Box, Prod_Wholesale_Min_Qty, Prod_Price_Wholesale_Unit, Prod_Sell_Only_Unit, Prod_Path, Prod_CategoryId, Prod_Stock FROM products"
	sentenciaCount := "SELECT count(*) FROM products"

	var join, where string
	args := []any{}

	switch choice {
	case "P":
		where = " WHERE Prod_Id = ?"
		args = append(args, p.ProdId)
	case "S":
		where = " WHERE UCASE(CONCAT_WS(' ', Prod_Title, Prod_Description)) LIKE ?"
		args = append(args, "%"+strings.ToUpper(p.ProdSearch)+"%")
	case "C":
		where = " WHERE Prod_CategoryId = ?"
		args = append(args, p.ProdCategId)
	case "U":
		where = " WHERE UCASE(Prod_Path) LIKE ?"
		args = append(args, "%"+strings.ToUpper(p.ProdPath)+"%")
	case "K":
		join = " JOIN category ON Prod_CategoryId = Categ_Id AND Categ_Path LIKE ?"
		args = append(args, "%"+strings.ToUpper(p.ProdCategPath)+"%")
	}

	var registros int
	err = Db.QueryRow(sentenciaCount+join+where, args...).Scan(&registros)
	if err != nil {
		fmt.Println(err.Error())
		return Resp, err
	}

	// orderField ya viene validado por el router; aquí solo se eligen textos fijos
	orderBy := ""
	switch orderField {
	case "I":
		orderBy = " ORDER BY Prod_Id"
	case "T":
		orderBy = " ORDER BY Prod_Title"
	case "D":
		orderBy = " ORDER BY Prod_Description"
	case "F":
		orderBy = " ORDER BY Prod_CreatedAt"
	case "P":
		orderBy = " ORDER BY Prod_Price"
	case "S":
		orderBy = " ORDER BY Prod_Stock"
	case "C":
		orderBy = " ORDER BY Prod_CategoryId"
	}
	if orderBy != "" && orderType == "D" {
		orderBy += " DESC"
	}
	if orderBy == "" {
		// Sin un orden fijo, la paginación puede repetir o saltarse productos
		orderBy = " ORDER BY Prod_Id"
	}

	limit := ""
	if page > 0 {
		if pageSize <= 0 {
			pageSize = 10
		}
		limit = " LIMIT ? OFFSET ?"
		args = append(args, pageSize, pageSize*(page-1))
	}

	rows, err := Db.Query(sentencia+join+where+orderBy+limit, args...)
	if err != nil {
		fmt.Println(err.Error())
		return Resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var prod models.Product
		var ProdId sql.NullInt32
		var ProdTitle sql.NullString
		var ProdDescription sql.NullString
		var ProdCreatedAt sql.NullTime
		var ProdUpdated sql.NullTime
		var ProdPrice sql.NullFloat64
		var ProdPriceBox sql.NullFloat64
		var ProdUnitsPerBox sql.NullInt32
		var ProdWholesaleMinQty sql.NullInt32
		var ProdPriceWholesaleUnit sql.NullFloat64
		var ProdSellOnlyUnit sql.NullInt32
		var ProdPath sql.NullString
		var ProdCategoryId sql.NullInt32
		var ProdStock sql.NullInt32

		err := rows.Scan(&ProdId, &ProdTitle, &ProdDescription, &ProdCreatedAt, &ProdUpdated, &ProdPrice, &ProdPriceBox, &ProdUnitsPerBox, &ProdWholesaleMinQty, &ProdPriceWholesaleUnit, &ProdSellOnlyUnit, &ProdPath, &ProdCategoryId, &ProdStock)
		if err != nil {
			return Resp, err
		}

		prod.ProdId = int(ProdId.Int32)
		prod.ProdTitle = ProdTitle.String
		prod.ProdDescription = ProdDescription.String
		prod.ProdCreatedAt = ProdCreatedAt.Time.String()
		prod.ProdUpdated = ProdUpdated.Time.String()
		prod.ProdPrice = ProdPrice.Float64
		prod.ProdPriceBox = ProdPriceBox.Float64
		prod.ProdUnitsPerBox = int(ProdUnitsPerBox.Int32)
		prod.ProdWholesaleMinQty = int(ProdWholesaleMinQty.Int32)
		prod.ProdPriceWholesaleUnit = ProdPriceWholesaleUnit.Float64
		prod.ProdSellOnlyUnit = int(ProdSellOnlyUnit.Int32)
		prod.ProdPath = ProdPath.String
		prod.ProdCategId = int(ProdCategoryId.Int32)
		prod.ProdStock = int(ProdStock.Int32)
		Prod = append(Prod, prod)
	}

	if err := rows.Err(); err != nil {
		return Resp, err
	}

	Resp.TotalItems = registros
	Resp.Data = Prod

	fmt.Println("Select Product > Ejecución Exitosa")
	return Resp, nil
}

func UpdateStock(p models.Product) error {
	fmt.Println("Comienza Update Stock")

	if p.ProdStock == 0 {
		return errors.New("[ERROR] Debe enviar el Stock a modificar")
	}

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	_, err = Db.Exec("UPDATE products SET Prod_Stock = Prod_Stock + ? WHERE Prod_Id = ?", p.ProdStock, p.ProdId)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Stock > Ejecución Exitosa")
	return nil
}