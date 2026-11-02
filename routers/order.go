package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/javier/santaMaria/bd"
	"github.com/javier/santaMaria/models"
)

// func InsertOrder(body string, User string) (int, string) {
// 	var o models.Orders

// 	err := json.Unmarshal([]byte(body), &o)
// 	if err != nil {
// 		return 400, "Error en los datos recibidos " + err.Error()
// 	}

// 	o.Order_UserUUID = User

// 	OK, message := ValidOrder(o)
// 	if !OK {
// 		return 400, message
// 	}

// 	result, err2 := bd.InsertOrder(o)
// 	if err2 != nil {
// 		return 400, "Ocurrió un error al intentar realizar el registro de la orden " + err2.Error()
// 	}

// 	return 200, "{ OrderID: " + strconv.Itoa(int(result)) + "}"
// }

// func ValidOrder(o models.Orders) (bool, string) {
// 	if o.Order_Total == 0 {
// 		return false, "Debe indicar el total de la orden"
// 	}

// 	count := 0
// 	for _, od := range o.OrderDetails {
// 		if od.OD_ProdId == 0{
// 			return false, "Debe indicar el ID del producto en el detalle de la Orden"
// 		}
// 		if od.OD_Quantity == 0{
// 			return false, "Debe indicar la cantidad del producto en el detalle de la Orden"
// 		}
// 		count++
// 	}
// 	if count == 0 {
// 		return false, "Debe indicar items en la orden"
// 	}
// 	return true, ""
// }

func InsertOrder(body string, User string) (int, string) {
	var req models.OrderRequest

	err := json.Unmarshal([]byte(body), &req)
	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	}

	if req.OrderAddId <= 0 {
		return 400, "Debe indicar la dirección de envío"
	}
	if len(req.OrderDetails) == 0 {
		return 400, "Debe indicar items en la orden"
	}
	if len(req.OrderDetails) > 100 {
		return 400, "La orden tiene demasiados items"
	}
	for _, od := range req.OrderDetails {
		if od.OD_ProdId <= 0 {
			return 400, "Debe indicar el ID del producto en el detalle de la Orden"
		}
		if od.OD_Quantity <= 0 {
			return 400, "La cantidad del producto en el detalle de la Orden debe ser mayor a 0"
		}
	}

	orderId, total, err := bd.InsertOrder(User, req)
	if err != nil {
		if errors.Is(err, bd.ErrInvalidOrder) {
			return 400, err.Error()
		}
		fmt.Println("InsertOrder error: " + err.Error())
		return 500, "Ocurrió un error al intentar registrar la orden"
	}

	return 200, fmt.Sprintf(`{"orderId":%d,"orderTotal":%.2f}`, orderId, total)
}

func SelectOrders(user string, request events.APIGatewayV2HTTPRequest) (int, string) {
	param := request.QueryStringParameters

	var page, orderId int
	var err error

	if len(param["page"]) > 0 {
		page, err = strconv.Atoi(param["page"])
		if err != nil || page < 1 {
			return 400, "El valor de page debe ser un número entero mayor a 0"
		}
	}
	if len(param["orderId"]) > 0 {
		orderId, err = strconv.Atoi(param["orderId"])
		if err != nil || orderId < 1 {
			return 400, "El valor de orderId debe ser un número entero mayor a 0"
		}
	}

	result, err := bd.SelectOrders(user, param["fechaDesde"], param["fechaHasta"], page, orderId)
	if err != nil {
		fmt.Println("SelectOrders error: " + err.Error())
		return 500, "Ocurrió un error al intentar obtener las órdenes"
	}

	Orders, err := json.Marshal(result)
	if err != nil {
		return 500, "Ocurrió un error al convertir las órdenes a JSON"
	}

	return 200, string(Orders)
}

// func SelectOrders(user string, request events.APIGatewayV2HTTPRequest) (int, string) {
// 	var fechaDesde, fechaHasta string
// 	var orderId int
// 	var page int

// 	if len(request.QueryStringParameters["fechaDesde"]) > 0 {
// 		fechaDesde = request.QueryStringParameters["fechaDesde"]
// 	}
// 	if len(request.QueryStringParameters["fechaHasta"]) > 0 {
// 		fechaHasta = request.QueryStringParameters["fechaHasta"]
// 	}
// 	if len(request.QueryStringParameters["page"]) > 0 {
// 		page, _ = strconv.Atoi(request.QueryStringParameters["page"])
// 	}
// 	if len(request.QueryStringParameters["orderId"]) > 0 {
// 		orderId, _ = strconv.Atoi(request.QueryStringParameters["orderId"])
// 	}
	
	
// 	result, err := bd.SelectOrders(user, fechaDesde,fechaHasta, page, orderId)
// 	if err != nil {
// 		return 400, "Ocurrió un error al intentar capturar los registros de órdenes del " + fechaDesde + " al " + fechaHasta + " > " + err.Error()
// 	}

// 	Orders, err2 := json.Marshal(result)
// 	if err2 != nil {
// 		return 400, "Ocurrió un error al intentar convertir en JSON el registro de Orden"
// 	}
	
// 	return 200, string(Orders)
// }