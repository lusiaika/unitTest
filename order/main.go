package main

import (
	"fmt"
	"log"
	"net/http"
	"order/database"
	"order/handler"
	"time"

	"github.com/gorilla/mux"
)

var PORT = ":8888"

// Replace with your own connection parameters
var sqlserver = "localhost"
var sqlport = 1433
var sqldbName = "orders_by"
var sqluser = "lusiaika"
var sqlpassword = "123456"

func main() {
	// Create connection string
	// connString := fmt.Sprintf("server=%s;database=%s;port=%d;trusted_connection=yes",
	// 	sqlserver, sqldbName, sqlport)

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		sqluser, sqlpassword, sqlserver, sqlport, sqldbName)

	sql := database.NewSqlConnection(connString)
	handler.SqlConnect = sql
	r := mux.NewRouter()

	userHandler := handler.NewOrderHandler()
	r.HandleFunc("/orders", userHandler.OrdersHandler)
	r.HandleFunc("/orders/{orderId}", userHandler.OrdersHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//http.ListenAndServe(PORT, nil)
	log.Fatal(srv.ListenAndServe())
}
