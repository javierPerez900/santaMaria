package routers

import (
	"encoding/json"
	"fmt"

	"github.com/javier/santaMaria/bd"
	"github.com/javier/santaMaria/models"
)

func InsertAddress(body string, User string) (int, string) {
	var t models.Address
	err := json.Unmarshal([]byte(body), &t)

	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	}

	if t.AddAddress == "" {
		return 400, "Debe especificar el Address"
	}
	if t.AddName == "" {
		return 400, "Debe especificar el Name"
	}
	if t.AddTitle == "" {
		return 400, "Debe especificar el Title"
	}
	if t.AddCity == "" {
		return 400, "Debe especificar el City"
	}
	if t.AddState == "" {
		return 400, "Debe especificar el State"
	}
	if t.AddPhone == "" {
		return 400, "Debe especificar el Phone"
	}
	if t.AddPostalCode == "" {
		return 400, "Debe especificar el PostalCode"
	}

	err = bd.InsertAddress(t, User)
	if err != nil {
		fmt.Println("InsertAddress error: " + err.Error())
		return 500, "Ocurrió un error al intentar registrar la dirección"
	}

	return 200, "InsertAddress OK"
}

func UpdateAddress(body string, User string, id int) (int, string){
	var t models.Address

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error en los datos recibidos" + err.Error()
	} 

	if t.AddAddress == "" && t.AddCity == "" && t.AddName == "" && t.AddPhone == "" &&
		t.AddPostalCode == "" && t.AddState == "" && t.AddTitle == "" {
		return 400, "Debe especificar al menos un campo para actualizar"
	}

	t.AddId = id
	err, encontrado := bd.AddressExists(User, t.AddId)
	if err != nil {
		fmt.Println("AddressExists error: " + err.Error())
		return 500, "Ocurrió un error al intentar buscar la dirección"
	}
	if !encontrado {
		return 404, "No se encuentra una dirección con ese ID para este usuario"
	}

	err = bd.UpdateAddress(t, User)
	if err != nil {
		fmt.Println("UpdateAddress error: " + err.Error())
		return 500, "Ocurrió un error al intentar actualizar la dirección"
	}

	return 200, "UpdateAddress OK"
}

func DeleteAddress(User string, id int) (int, string){

	err, encontrado := bd.AddressExists(User, id)
	if err != nil {
		fmt.Println("AddressExists error: " + err.Error())
		return 500, "Ocurrió un error al intentar buscar la dirección"
	}
	if !encontrado {
		return 404, "No se encuentra una dirección con ese ID para este usuario"
	}

	err = bd.DeleteAddress(id, User)
	if err != nil {
		fmt.Println("DeleteAddress error: " + err.Error())
		return 500, "Ocurrió un error al intentar borrar la dirección"
	}

	return 200, "DeleteAddress OK"
}

func SelectAddress(User string) (int, string) {
	addr, err := bd.SelectAddress(User)
	if err != nil {
		fmt.Println("SelectAddress error: " + err.Error())
		return 500, "Ocurrió un error al intentar obtener las direcciones"
	}

	respJson, err := json.Marshal(addr)
	if err != nil {
		return 500, "Ocurrió un error al formatear las direcciones como JSON"
	}

	return 200, string(respJson)
}