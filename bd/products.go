package bd

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/javier/santaMaria/models"
	"github.com/javier/santaMaria/tools"
)

func InsertProduct(p models.Product) (int64,error) {
	fmt.Println("Comienza Registro")

	err := DbConnect()
	if err != nil {
		return 0, err
	}
	defer Db.Close()

	sentencia := "INSERT INTO products (Prod_Title "

	if len(p.ProdDescription)>0 {
		sentencia += ", Prod_Description"
	}
	if p.ProdPrice > 0 {
		sentencia += ", Prod_Price"
	}
	if p.ProdCategId > 0 {
		sentencia += ", Prod_CategoryId"
	}
	if p.ProdStock > 0 {
		sentencia += ", Prod_Stock"
	}
	if len(p.ProdPath) > 0 {
		sentencia += ", Prod_Path"
	}
	if p.ProdSellOnlyUnit == 0 || p.ProdSellOnlyUnit == 1 {
		sentencia += ", Prod_Sell_Only_Unit"

		if p.ProdSellOnlyUnit == 0 {
			if p.ProdPriceBox > 0 {
				sentencia += ", Prod_Price_Box"
			}
			if p.ProdUnitsPerBox > 0 {
				sentencia += ", Prod_Units_Per_Box"
			}
			if p.ProdWholesaleMinQty > 0 {
				sentencia += ", Prod_Wholesale_Min_Qty"
			}
			if p.ProdPriceBox > 0 && p.ProdUnitsPerBox > 0{
				sentencia += ", Prod_Price_Wholesale_Unit"
			}
		}

	}

	sentencia += ") VALUES ('" + tools.EscapeString(p.ProdTitle) + "'"

	if len(p.ProdDescription)>0 {
		sentencia += ", '" + tools.EscapeString(p.ProdDescription) + "'"
	}
	if p.ProdPrice > 0 {	
		sentencia += ", " + strconv.FormatFloat(p.ProdPrice, 'e', -1, 64)
	}
	if p.ProdCategId > 0 {
		sentencia += ", " + strconv.Itoa(p.ProdCategId)
	}
	if p.ProdStock > 0 {
		sentencia += ", " + strconv.Itoa(p.ProdStock)
	}
	if len(p.ProdPath) > 0 {
		sentencia += ", '" + tools.EscapeString(p.ProdPath) + "'"
	}
	if p.ProdSellOnlyUnit == 0 || p.ProdSellOnlyUnit == 1 {
		sentencia += ", " + strconv.Itoa(p.ProdSellOnlyUnit)

		if p.ProdSellOnlyUnit == 0 {
			if p.ProdPriceBox > 0 {
				sentencia += ", " + strconv.FormatFloat(p.ProdPriceBox, 'e', -1, 64)
			}
			if p.ProdUnitsPerBox > 0 {
				sentencia += ", " + strconv.Itoa(p.ProdUnitsPerBox)
			}
			if p.ProdWholesaleMinQty > 0 {
				sentencia += ", " + strconv.Itoa(p.ProdWholesaleMinQty)
			}
			if p.ProdPriceBox > 0 && p.ProdUnitsPerBox > 0{
			priceWholesaleUnit := p.ProdPriceBox / float64(p.ProdUnitsPerBox)
				sentencia += ", " + strconv.FormatFloat(priceWholesaleUnit, 'e', -1, 64)
			}
		}

	}

	sentencia += ")"

	var result sql.Result
	result, err = Db.Exec(sentencia)
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

	err := DbConnect()
	if err!=nil {
		return err
	}
	defer Db.Close()

	sentencia := "UPDATE products SET "

	sentencia = tools.ArmoSentencia(sentencia, "Prod_Title", "S", 0, 0, p.ProdTitle)
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Description", "S", 0, 0, p.ProdDescription)
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Price", "F", 0, p.ProdPrice, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_CategoryId", "N", p.ProdCategId, 0, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Stock", "N", p.ProdStock, 0, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Path", "S", 0, 0, p.ProdPath)
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Sell_Only_Unit", "N", p.ProdSellOnlyUnit, 0, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Price_Box", "F", 0, p.ProdPriceBox, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Units_Per_Box", "N", p.ProdUnitsPerBox, 0, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Wholesale_Min_Qty", "N", p.ProdWholesaleMinQty, 0, "")
	if p.ProdPriceBox > 0 && p.ProdUnitsPerBox > 0{
		priceWholesaleUnit := p.ProdPriceBox / float64(p.ProdUnitsPerBox)
		sentencia = tools.ArmoSentencia(sentencia, "Prod_Price_Wholesale_Unit", "F", 0, priceWholesaleUnit, "")
	}
	
	sentencia += " WHERE Prod_id = " + strconv.Itoa(p.ProdId)

	_, err = Db.Exec(sentencia)
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
	if err!=nil {
		return err
	}
	defer Db.Close()

	sentencia := "DELETE FROM products WHERE Prod_Id = " + strconv.Itoa(id)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(sentencia)
	fmt.Println("Delete Product > Ejecución Exitosa")
	return nil
}

func SelectProduct(p models.Product, choice string, page int, pageSize int, orderType string, orderField string) (models.ProductResp, error) {
	fmt.Println("Comienza SelectProduct")
	var Resp models.ProductResp
	var Prod []models.Product

	err := DbConnect()
	if err!=nil {
		return Resp, err
	}
	defer Db.Close()

	var sentencia string
	var sentenciaCount string
	var where, limit string

	sentencia = "SELECT Prod_Id, Prod_Title, Prod_Description, Prod_CreatedAt, Prod_Updated, Prod_Price, Prod_Price_Box, Prod_Units_Per_Box, Prod_Wholesale_Min_Qty, Prod_Price_Wholesale_Unit, Prod_Sell_Only_Unit, Prod_Path, Prod_CategoryId, Prod_Stock FROM products"
	sentenciaCount = "SELECT count(*) as registros FROM products"

	switch choice {
	case "P":
		where = " WHERE Prod_Id = " + strconv.Itoa(p.ProdId)
	case "S":
		where = " WHERE UCASE(CONCAT(Prod_Title, Prod_Description)) LIKE '%" + strings.ToUpper(p.ProdSearch) + "%' "
	case "C":
		where = " WHERE Prod_CategoryId = " + strconv.Itoa(p.ProdCategId)
	case "U":
		where = " WHERE UCASE(Prod_Path) LIKE '%" + strings.ToUpper(p.ProdPath) + "%' "
	case "K":
		join := " Join category ON Prod_CategoryId = Categ_Id AND Categ_Path LIKE '%" + strings.ToUpper(p.ProdCategPath) + "%' "
		sentencia += join
		sentenciaCount += join
	}

	sentenciaCount += where

	var rows *sql.Rows
	rows, err = Db.Query(sentenciaCount)

	if err != nil {
		fmt.Println(err.Error())
		return Resp, err
	}

	rows.Next() 
	var regi sql.NullInt32
	err = rows.Scan(&regi)

	registros := int(regi.Int32)

	if page > 0 {
		if registros > pageSize {
			limit = " LIMIT " + strconv.Itoa(pageSize)
			if page > 1 {
				offset := pageSize * ( page - 1 )
				limit += " OFFSET " + strconv.Itoa(offset)
			}
		} else {
			limit = ""
		}
	}

	var orderBy string
	if len(orderField) > 0 {
		switch orderField  {
		case "I":
			orderBy = " ORDER BY Prod_Id "

		case "T":
			orderBy = " ORDER BY Prod_Title "
		
		case "D":
			orderBy = " ORDER BY Prod_Description "
	
		case "F":
			orderBy = " ORDER BY Prod_CreatedAt "
		
		case "P":
			orderBy = " ORDER BY Prod_Price "
		
		case "S":
			orderBy = " ORDER BY Prod_Stock "
		
		case "C":
			orderBy = " ORDER BY Prod_CategoryId"
		}

		if orderType == "D" {
			orderBy	+= " DESC"
		}
	}

	sentencia += where + orderBy + limit

	fmt.Println(sentencia)

	rows, err = Db.Query(sentencia)

	for rows.Next() {
		var p models.Product
		var ProdId sql.NullInt32
		var ProdTitle sql.NullString
		var ProdDescription sql.NullString
		var ProdCreatedAt sql.NullTime
		var ProdUpdated sql.NullTime
		var ProdPrice sql.NullFloat64
		var ProdPriceBox sql.NullFloat64
		var ProdUnitsPerBox sql.NullInt32
		var ProdWholesaleMinQty  sql.NullInt32
		var ProdPriceWholesaleUnit sql.NullFloat64
		var ProdSellOnlyUnit sql.NullInt32
		var ProdPath sql.NullString
		var ProdCategoryId sql.NullInt32
		var ProdStock sql.NullInt32
	
		err := rows.Scan(&ProdId, &ProdTitle, &ProdDescription, &ProdCreatedAt, &ProdUpdated, &ProdPrice, &ProdPriceBox, &ProdUnitsPerBox, &ProdWholesaleMinQty, &ProdPriceWholesaleUnit, &ProdSellOnlyUnit, &ProdPath, &ProdCategoryId, &ProdStock)
		if err != nil {
			return Resp, err
		}

		p.ProdId = int(ProdId.Int32)
		p.ProdTitle = ProdTitle.String
		p.ProdDescription = ProdDescription.String
		p.ProdCreatedAt = ProdCreatedAt.Time.String()
		p.ProdUpdated = ProdUpdated.Time.String()
		p.ProdPrice = ProdPrice.Float64
		p.ProdPriceBox = ProdPriceBox.Float64
		p.ProdUnitsPerBox = int(ProdUnitsPerBox.Int32)
		p.ProdWholesaleMinQty = int(ProdWholesaleMinQty.Int32)
		p.ProdPriceWholesaleUnit = ProdPriceWholesaleUnit.Float64
		p.ProdSellOnlyUnit = int(ProdSellOnlyUnit.Int32)
		p.ProdPath = ProdPath.String
		p.ProdCategId = int(ProdCategoryId.Int32)
		p.ProdStock = int(ProdStock.Int32)
		Prod = append(Prod, p)
	}

	Resp.TotalItems = registros
	Resp.Data = Prod

	fmt.Println("Select Product > Ejecución Exitosa")
	return Resp, nil
}

func UpdateStock(p models.Product) error {
	fmt.Println("Comienza Update Stock")

	if p.ProdStock==0 {
		return errors.New("[ERROR] Debe enviar el Stock a modificar")
	}

	err := DbConnect()
	if err!=nil {
		return err
	}
	defer Db.Close()

	sentencia := "UPDATE products SET Prod_Stock = Prod_Stock + " + strconv.Itoa(p.ProdStock) + " WHERE Prod_Id = " + strconv.Itoa(p.ProdId)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Stock > Ejecución Exitosa")
	return nil
}