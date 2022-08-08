package entity

import "time"

type Item struct {
	ItemID      int    `json:"lineItemId"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type Orders struct {
	OrderID      int       `json:"orderId"`
	OrderedAt    time.Time `json:"orderedAt"`
	CustomerName string    `json:"customerName"`
	Items        []Item    `json:"items,omitempty"`
}
type DataItem struct {
	ItemCode    string `tvp:"item_code"`
	Description string `tvp:"description"`
	Quantity    int    `tvp:"quantity"`
}
