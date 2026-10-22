package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/javier/santaMaria/bd"
	"github.com/javier/santaMaria/models"
)

func UpdateUser(body string, User string) (int, string) {
	var t models.User
	err := json.Unmarshal([]byte(body), &t)

	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	}

	if len(t.UserFirstName) == 0 && len(t.UserLastName) == 0 {
		return 400, "Debe especificar el Nombre (FirstName) o (LastName) del usuario "
	}

	err, encontrado := bd.UserExists(User)
	if err != nil {
		fmt.Println("UserExists error: " + err.Error())
		return 500, "Ocurrió un error al intentar buscar el usuario"
	}
	if !encontrado {
		return 404, "No existe un usuario con ese UUID '" + User + "'"
	}

	err = bd.UpdateUser(t, User)
	if err != nil {
		fmt.Println("UpdateUser error: " + err.Error())
		return 500, "Ocurrió un error al intentar actualizar el usuario"
	}

	return 200, "UpdateUser OK"
}

func SelectUser(body string, User string) (int, string) {
	err, encontrado := bd.UserExists(User)
	if err != nil {
		fmt.Println("UserExists error: " + err.Error())
		return 500, "Ocurrió un error al intentar buscar el usuario"
	}
	if !encontrado {
		return 404, "No existe un usuario con ese UUID '" + User + "'"
	}

	row, err := bd.SelectUser(User)
	if err != nil {
		fmt.Println("SelectUser error: " + err.Error())
		return 500, "Ocurrió un error al intentar obtener los datos del usuario"
	}

	respJson, err := json.Marshal(row)
	if err != nil {
		return 500, "Error al formatear los datos del usuario como JSON"
	}
	return 200, string(respJson)
}

func SelectUsers(body string, User string, request events.APIGatewayV2HTTPRequest) (int, string) {
	isAdmin, msg := bd.UserIsAdmin(User)
	if !isAdmin {
		fmt.Println("UserIsAdmin: " + msg)
		return 403, "No tiene permisos para realizar esta acción"
	}

	Page := 1
	if p := request.QueryStringParameters["page"]; len(p) > 0 {
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 {
			return 400, "El valor de page debe ser un número entero mayor a 0"
		}
		Page = n
	}

	user, err := bd.SelectUsers(Page)
	if err != nil {
		fmt.Println("SelectUsers error: " + err.Error())
		return 500, "Ocurrió un error al intentar obtener la lista de usuarios"
	}

	respJson, err := json.Marshal(user)
	if err != nil {
		return 500, "Error al formatear los datos de los usuarios como JSON"
	}

	return 200, string(respJson)
}