package handler

import (
	"encoding/json"
	"net/http"
	"order/database"
)

var SqlConnect *database.Database

type response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

const (
	statusSuccess int = 0
	statusError   int = 1
)

func writeJsonResp(w http.ResponseWriter, status int, obj interface{}) {

	resp := response{
		Status: status,
		Data:   obj,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
