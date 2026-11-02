package models

type OrderRequest struct {
	OrderAddId   int                `json:"orderAddId"`
	OrderDetails []OrderItemRequest `json:"orderDetails"`
}

type OrderItemRequest struct {
	OD_ProdId   int `json:"odProdId"`
	OD_Quantity int `json:"odQuantity"`
}