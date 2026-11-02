package bd

import (
	"database/sql"
	"errors"
	"fmt"
	"math"

	_ "github.com/go-sql-driver/mysql"
	"github.com/javier/santaMaria/models"
)

var ErrInvalidOrder = errors.New("pedido inválido")

// precioProducto reúne los datos que intervienen en el precio, en céntimos.
type precioProducto struct {
	precioCents     int64
	precioCajaCents int64
	unidadesPorCaja int64
	minPorMayor     int64
	soloUnidad      bool
}

func aCentimos(valor float64) int64 {
	return int64(math.Round(valor * 100))
}

// esPorMayor indica si la cantidad (en unidades) llega al precio por mayor.
// Un mínimo de 0 significa que el producto no tiene precio por mayor.
func (p precioProducto) esPorMayor(unidades int64) bool {
	return !p.soloUnidad &&
		p.minPorMayor > 0 &&
		unidades >= p.minPorMayor &&
		p.unidadesPorCaja > 0 &&
		p.precioCajaCents > 0
}

// totalCents calcula el total de una línea en céntimos:
//   - al por menor: precio unitario × unidades
//   - por mayor: precio de caja × unidades / unidades por caja
//     (equivale a cajas × precio de caja + sobrantes × precio de caja / unidades por caja)
func (p precioProducto) totalCents(unidades int64) int64 {
	if p.esPorMayor(unidades) {
		return (p.precioCajaCents*unidades + p.unidadesPorCaja/2) / p.unidadesPorCaja
	}
	return p.precioCents * unidades
}

type lineaOrden struct {
	prodId       int
	unidades     int64
	precioUnit   float64
	totalCents   int64
	presentacion string
}

// func InsertOrder(o models.Orders) (int64, error) {
// 	fmt.Println("Comienza Registro de InsertOrder")

// 	err := DbConnect()
// 	if err!=nil {
// 		return 0, err
// 	}
// 	defer Db.Close()

// 	sentencia := "INSERT INTO orders (Order_UserUUID, Order_Total, Order_AddId) VALUES ('" 
// 	sentencia += o.Order_UserUUID + "', '" + strconv.FormatFloat(o.Order_Total, 'f', -1, 64) + "', '" + strconv.Itoa(o.Order_AddId) + "')"

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

// 	for _, od := range o.OrderDetails {
// 		sentencia = "INSERT INTO orders_detail (OD_OrderId, OD_ProdId, OD_Quantity, OD_Presentation, OD_Unit_Price, OD_Price) VALUES (" + strconv.Itoa(int(LastInsertId))
// 		sentencia += ", " + strconv.Itoa(od.OD_ProdId) + ", " +strconv.Itoa(od.OD_Quantity) + ", '"+ od.OD_Presentation + "', " + strconv.FormatFloat(od.OD_Unit_Price, 'f', -1, 64) + ", " + strconv.FormatFloat(od.OD_Price, 'f', -1, 64) + ")"

// 		fmt.Println(sentencia)
// 		_, err = Db.Exec(sentencia)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return 0, err
// 		}
// 	}	
	
// 	fmt.Println("Insert Order > Ejecución Exitosa")
// 	return LastInsertId, nil
// }

func InsertOrder(userUUID string, req models.OrderRequest) (int64, float64, error) {
	fmt.Println("Comienza Registro de InsertOrder")

	err := DbConnect()
	if err != nil {
		return 0, 0, err
	}
	defer Db.Close()

	tx, err := Db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback() // no hace nada si ya se confirmó con Commit

	// La dirección tiene que ser del usuario que compra
	var existe int
	err = tx.QueryRow("SELECT 1 FROM addresses WHERE Add_Id = ? AND Add_UserId = ?", req.OrderAddId, userUUID).Scan(&existe)
	if err == sql.ErrNoRows {
		return 0, 0, fmt.Errorf("%w: la dirección indicada no existe", ErrInvalidOrder)
	}
	if err != nil {
		return 0, 0, err
	}

	// Une las líneas repetidas del mismo producto, conservando el orden
	unidadesPorProducto := map[int]int64{}
	ordenProductos := []int{}
	for _, item := range req.OrderDetails {
		if _, ok := unidadesPorProducto[item.OD_ProdId]; !ok {
			ordenProductos = append(ordenProductos, item.OD_ProdId)
		}
		unidadesPorProducto[item.OD_ProdId] += int64(item.OD_Quantity)
	}

	lineas := make([]lineaOrden, 0, len(ordenProductos))
	var totalCents int64

	for _, prodId := range ordenProductos {
		unidades := unidadesPorProducto[prodId]

		var precio, precioCaja sql.NullFloat64
		var unidadesPorCaja, minPorMayor, soloUnidad sql.NullInt64

		err = tx.QueryRow(
			"SELECT Prod_Price, Prod_Price_Box, Prod_Units_Per_Box, Prod_Wholesale_Min_Qty, Prod_Sell_Only_Unit FROM products WHERE Prod_Id = ?",
			prodId,
		).Scan(&precio, &precioCaja, &unidadesPorCaja, &minPorMayor, &soloUnidad)
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("%w: el producto %d no existe", ErrInvalidOrder, prodId)
		}
		if err != nil {
			return 0, 0, err
		}

		p := precioProducto{
			precioCents:     aCentimos(precio.Float64),
			precioCajaCents: aCentimos(precioCaja.Float64),
			unidadesPorCaja: unidadesPorCaja.Int64,
			minPorMayor:     minPorMayor.Int64,
			soloUnidad:      soloUnidad.Int64 == 1,
		}
		if p.precioCents <= 0 {
			return 0, 0, fmt.Errorf("%w: el producto %d no tiene precio", ErrInvalidOrder, prodId)
		}

		totalLinea := p.totalCents(unidades)
		presentacion := "Unidad"
		if p.esPorMayor(unidades) {
			presentacion = "Por mayor"
		}

		lineas = append(lineas, lineaOrden{
			prodId:       prodId,
			unidades:     unidades,
			precioUnit:   float64(totalLinea) / float64(unidades) / 100,
			totalCents:   totalLinea,
			presentacion: presentacion,
		})
		totalCents += totalLinea
	}

	total := float64(totalCents) / 100

	result, err := tx.Exec(
		"INSERT INTO orders (Order_UserUUID, Order_Total, Order_AddId) VALUES (?, ?, ?)",
		userUUID, total, req.OrderAddId,
	)
	if err != nil {
		return 0, 0, err
	}

	orderId, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	for _, l := range lineas {
		_, err = tx.Exec(
			"INSERT INTO orders_detail (OD_OrderId, OD_ProdId, OD_Quantity, OD_Presentation, OD_Unit_Price, OD_Price) VALUES (?, ?, ?, ?, ?, ?)",
			orderId, l.prodId, l.unidades, l.presentacion, l.precioUnit, float64(l.totalCents)/100,
		)
		if err != nil {
			return 0, 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, 0, err
	}

	fmt.Println("Insert Order > Ejecución Exitosa")
	return orderId, total, nil
}

func SelectOrders(user string, fechaDesde string, fechaHasta string, page int, orderId int) ([]models.Orders, error) {
	fmt.Println("Comienza SelectOrders")

	// Slice vacío (no nil): el JSON sale como [] y no como null cuando no hay pedidos
	Orders := []models.Orders{}

	err := DbConnect()
	if err != nil {
		return Orders, err
	}
	defer Db.Close()

	sentencia := "SELECT Order_Id, Order_UserUUID, Order_AddId, Order_Date, Order_Total FROM orders WHERE Order_UserUUID = ?"
	args := []any{user}

	if orderId > 0 {
		sentencia += " AND Order_Id = ?"
		args = append(args, orderId)
	} else {
		if len(fechaDesde) > 0 && len(fechaHasta) > 0 {
			if len(fechaHasta) == 10 {
				fechaHasta += " 23:59:59"
			}
			sentencia += " AND Order_Date BETWEEN ? AND ?"
			args = append(args, fechaDesde, fechaHasta)
		}

		if page < 1 {
			page = 1
		}
		sentencia += " ORDER BY Order_Id LIMIT 10 OFFSET ?"
		args = append(args, 10*(page-1))
	}

	rows, err := Db.Query(sentencia, args...)
	if err != nil {
		return Orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var Order models.Orders
		var OrderAddId sql.NullInt32

		err := rows.Scan(&Order.Order_Id, &Order.Order_UserUUID, &OrderAddId, &Order.Order_Date, &Order.Order_Total)
		if err != nil {
			return Orders, err
		}

		Order.Order_AddId = int(OrderAddId.Int32)
		Orders = append(Orders, Order)
	}
	if err := rows.Err(); err != nil {
		return Orders, err
	}
	rows.Close() 

	for i := range Orders {
		rowsD, err := Db.Query(
			"SELECT OD_Id, OD_ProdId, OD_Quantity, OD_Presentation, OD_Unit_Price, OD_Price FROM orders_detail WHERE OD_OrderId = ?",
			Orders[i].Order_Id,
		)
		if err != nil {
			return Orders, err
		}

		for rowsD.Next() {
			var odId, odProdId, odQuantity sql.NullInt64
			var odPresentation sql.NullString
			var odUnitPrice, odPrice sql.NullFloat64

			err := rowsD.Scan(&odId, &odProdId, &odQuantity, &odPresentation, &odUnitPrice, &odPrice)
			if err != nil {
				rowsD.Close()
				return Orders, err
			}

			Orders[i].OrderDetails = append(Orders[i].OrderDetails, models.OrdersDetails{
				OD_Id:           int(odId.Int64),
				OD_ProdId:       int(odProdId.Int64),
				OD_Quantity:     int(odQuantity.Int64),
				OD_Presentation: odPresentation.String,
				OD_Unit_Price:   odUnitPrice.Float64,
				OD_Price:        odPrice.Float64,
			})
		}

		if err := rowsD.Err(); err != nil {
			rowsD.Close()
			return Orders, err
		}
		rowsD.Close()
	}

	fmt.Println("Select Orders > Ejecución Exitosa")
	return Orders, nil
}