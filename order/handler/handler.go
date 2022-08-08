package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"order/entity"
	"strconv"

	"github.com/gorilla/mux"
)

type OrderHandlerInterface interface {
	OrdersHandler(w http.ResponseWriter, r *http.Request)
}

type OrderHandler struct {
	//postgrespool *pgxpool.Pool
}

func NewOrderHandler() OrderHandlerInterface {
	//return &UserHandler{postgrespool: postgrespool}
	return &OrderHandler{}
}

func (h *OrderHandler) OrdersHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["orderId"]

	switch r.Method {
	case http.MethodGet:
		if id != "" { // get by id
			getOrderByIDHandler(w, r, id)
		} else { // get all
			h.getOrdersHandler(w, r)
		}
	case http.MethodPost:
		createOrderHandler(w, r)
	case http.MethodPut:
		updateOrderHandler(w, r, id)
	case http.MethodDelete:
		deleteOrderHandler(w, r, id)
	}
}

func (h *OrderHandler) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	orders, err := SqlConnect.GetOrders(ctx)
	for i, order := range orders {
		if items, err := SqlConnect.GetItems(ctx, &order); err == nil {

			orders[i].Items = items
		}
	}
	if err != nil {
		writeJsonResp(w, statusError, err.Error())
		return
	}
	writeJsonResp(w, statusSuccess, orders)
}

func getOrderByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	if idInt, err := strconv.Atoi(id); err == nil {
		ctx := context.Background()
		order, err := SqlConnect.GetOrderByID(ctx, idInt)
		if err != nil {
			writeJsonResp(w, statusError, err.Error())
			return
		} else {
			if items, err := SqlConnect.GetItems(ctx, order); err == nil {

				order.Items = items
			}
		}
		writeJsonResp(w, statusSuccess, order)
	}
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	decoder := json.NewDecoder(r.Body)
	var order entity.Orders
	if err := decoder.Decode(&order); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}

	orders, err := SqlConnect.CreateOrder(ctx, order)
	if err != nil {
		writeJsonResp(w, statusError, err.Error())
		return
	}
	writeJsonResp(w, statusSuccess, orders)
}

func updateOrderHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()

	if id != "" { // get by id
		decoder := json.NewDecoder(r.Body)
		var order entity.Orders
		if err := decoder.Decode(&order); err != nil {
			w.Write([]byte("error decoding json body"))
			return
		}

		if idInt, err := strconv.Atoi(id); err == nil {
			if orders, err := SqlConnect.GetOrderByID(ctx, idInt); err != nil {
				writeJsonResp(w, statusError, err.Error())
				return
			} else if orders.OrderID == 0 {
				writeJsonResp(w, statusError, "Data not exists")
				return
			} else {
				result, err := SqlConnect.UpdateOrder(ctx, idInt, order)
				if err != nil {
					writeJsonResp(w, statusError, err.Error())
					return
				}
				writeJsonResp(w, statusSuccess, result)
			}
		}
	}
}

func deleteOrderHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if id != "" { // get by id
		if idInt, err := strconv.Atoi(id); err == nil {
			result, err := SqlConnect.DeleteOrder(ctx, idInt)
			if err != nil {
				writeJsonResp(w, statusError, err.Error())
				return
			}
			writeJsonResp(w, statusSuccess, result)
		}
	}
}
