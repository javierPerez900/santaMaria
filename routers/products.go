package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/javier/santaMaria/bd"
	"github.com/javier/santaMaria/models"
)

func InsertProduct(body string, User string) (int, string) {
	var t models.Product
	err := json.Unmarshal([]byte(body), &t)

	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	}

	if len(t.ProdTitle) == 0 {
		return 400, "Debe especificar el Nombre (Title) del producto"
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		fmt.Println("UserIsAdmin: " + msg)
		return 403, "No tiene permisos para realizar esta acción"
	}

	result, err2 := bd.InsertProduct(t)
	if err2 != nil {
		fmt.Println("InsertProduct error: " + err2.Error())
		return 500, "Ocurrió un error al intentar registrar el producto"
	}

	return 200, fmt.Sprintf(`{"prodID":%d}`, result)
}

func UpdateProduct(body string, User string, id int) (int, string){
	var t models.Product

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		fmt.Println("UserIsAdmin: " + msg)
		return 403, "No tiene permisos para realizar esta acción"
	}

	t.ProdId = id
	err2 := bd.UpdateProduct(t)
	if err2 != nil {
		if errors.Is(err2, bd.ErrNoFields) {
			return 400, "Debe especificar al menos un campo para actualizar"
		}
		fmt.Println("UpdateProduct error: " + err2.Error())
		return 500, "Ocurrió un error al intentar actualizar el producto"
	}

	return 200, "Update OK"
}

func DeleteProduct(User string, id int) (int, string) {
	if id == 0 {
		return 400, "Debe especificar ID del producto a borrar"
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		fmt.Println("UserIsAdmin: " + msg)
		return 403, "No tiene permisos para realizar esta acción"
	}

	err := bd.DeleteProduct(id)
	if err != nil {
		fmt.Println("DeleteProduct error: " + err.Error())
		return 500, "Ocurrió un error al intentar borrar el producto"
	}

	return 200, "Delete OK"
}

func SelectProduct(request events.APIGatewayV2HTTPRequest) (int, string) {
	var t models.Product
	var page, pageSize int
	var err error

	param := request.QueryStringParameters

	if len(param["page"]) > 0 {
		page, err = strconv.Atoi(param["page"])
		if err != nil || page < 1 {
			return 400, "El valor de page debe ser un número entero mayor a 0"
		}
	}
	if len(param["pageSize"]) > 0 {
		pageSize, err = strconv.Atoi(param["pageSize"])
		if err != nil || pageSize < 1 || pageSize > 200 {
			return 400, "El valor de pageSize debe ser un número entero entre 1 y 200"
		}
	}

	orderType := param["orderType"]   // D = Desc. A o vacío = Asc
	orderField := param["orderField"] // 'I' Id, 'T' Title, 'D' Description, 'F' Created At
	//                                   'P' Price, 'C' CategId, 'S' Stock

	if len(orderField) != 1 || !strings.Contains("ITDFPCS", orderField) {
		orderField = ""
	}

	var choice string

	if len(param["prodId"]) > 0 {
		choice = "P"
		t.ProdId, err = strconv.Atoi(param["prodId"])
		if err != nil {
			return 400, "El valor de prodId debe ser un número entero"
		}
	}
	if len(param["search"]) > 0 {
		choice = "S"
		t.ProdSearch = param["search"]
	}
	if len(param["categId"]) > 0 {
		choice = "C"
		t.ProdCategId, err = strconv.Atoi(param["categId"])
		if err != nil {
			return 400, "El valor de categId debe ser un número entero"
		}
	}
	if len(param["slug"]) > 0 {
		choice = "U"
		t.ProdPath = param["slug"]
	}
	if len(param["slugCateg"]) > 0 {
		choice = "K"
		t.ProdCategPath = param["slugCateg"]
	}

	result, err2 := bd.SelectProduct(t, choice, page, pageSize, orderType, orderField)
	if err2 != nil {
		fmt.Println("SelectProduct error: " + err2.Error())
		return 500, "Ocurrió un error al intentar obtener los productos"
	}

	Product, err3 := json.Marshal(result)
	if err3 != nil {
		return 500, "Ocurrió un error al convertir los productos a JSON"
	}

	return 200, string(Product)
}

func UpdateStock(body string, User string, id int) (int, string){
	var t models.Product

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	}

	if t.ProdStock == 0 {
		return 400, "Debe enviar el Stock a modificar"
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		fmt.Println("UserIsAdmin: " + msg)
		return 403, "No tiene permisos para realizar esta acción"
	}

	t.ProdId = id
	err2 := bd.UpdateStock(t)
	if err2 != nil {
		fmt.Println("UpdateStock error: " + err2.Error())
		return 500, "Ocurrió un error al intentar actualizar el stock"
	}

	return 200, "Update OK"
}