package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/javier/santaMaria/bd"
	"github.com/javier/santaMaria/models"
)

func InsertCategory(body string, User string) (int, string) {
	var t models.Category

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	}

	if len(t.CategName) == 0 {
		return 400, "Debe especificar el Nombre (Title) de la categoría"
	}

	if len(t.CategPath) == 0 {
		return 400, "Debe especificar el Path (Ruta) de la categoría"
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		fmt.Println("UserIsAdmin: " + msg)
		return 403, "No tiene permisos para realizar esta acción"
	}

	result, err2 := bd.InsertCategory(t)
	if err2 != nil {
		fmt.Println("InsertCategory error: " + err2.Error())
		return 500, "Ocurrió un error al intentar registrar la categoría"
	}

	return 200, fmt.Sprintf(`{"categID":%d}`, result)
}

func UpdateCategory(body string, User string, id int) (int, string){
	var t models.Category

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	} 

	if len(t.CategName) == 0 && len(t.CategPath) == 0 {
		return 400, "Debe especificar CategName y CategPath para actualizar"
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		fmt.Println("UserIsAdmin: " + msg)
		return 403, "No tiene permisos para realizar esta acción"
	}

	t.CategID=id
	err2 := bd.UpdateCategory(t)
	if err2 != nil {
		fmt.Println("UpdateCategory error: " + err2.Error())
		return 500, "Ocurrió un error al intentar actualizar la categoría"
	}

	return 200, "Update OK"
}

func DeleteCategory(body string, User string, id int) (int, string) {
	if id == 0 {
		return 400, "Debe especificar ID de la categoría a borrar"
	}

	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		fmt.Println("UserIsAdmin: " + msg)
		return 403, "No tiene permisos para realizar esta acción"
	}

	err := bd.DeleteCategory(id)
	if err != nil {
		fmt.Println("DeleteCategory error: " + err.Error())
		return 500, "Ocurrió un error al intentar borrar la categoría"
	}

	return 200, "Delete OK"
}

func SelectCategories(body string, request events.APIGatewayV2HTTPRequest) (int, string) {
	var err error
	var CategId int
	var Slug string

	if len(request.QueryStringParameters["categId"]) > 0 {
		CategId, err = strconv.Atoi(request.QueryStringParameters["categId"])
		if err != nil {
			return 400, "El valor de categId debe ser un número entero"
		}
	} else {
		if len(request.QueryStringParameters["slug"]) > 0 {
			Slug = request.QueryStringParameters["slug"]
		}
	}

	lista, err2 := bd.SelectCategories(CategId, Slug)
	if err2 != nil {
		fmt.Println("SelectCategories error: " + err2.Error())
		return 500, "Ocurrió un error al intentar obtener las categorías"
	}

	Categ, err3 := json.Marshal(lista)
	if err3 != nil {
		return 500, "Ocurrió un error al convertir las categorías a JSON"
	}
	
	return 200, string(Categ)
}