package handler

import (
	"autoreloaddata/entity"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

const htmlPath = "html/web.html"
const jsonPath = "html/status.json"

func StatusHandler(w http.ResponseWriter, req *http.Request) {
	var data entity.DataStatus
	w.Header().Add("Content Type", "text/html")
	// read from json file and write to webData
	file, _ := ioutil.ReadFile(jsonPath)
	json.Unmarshal(file, &data)
	template, _ := template.ParseFiles(htmlPath)
	dt := CheckStatus(data)
	context := dt
	template.Execute(w, context)
}

func CheckStatus(s entity.DataStatus) *entity.DataStatus {
	if s.Status.Water <= 5 || s.Status.Wind <= 6 {
		s.Status.DataStatus = "Aman"
	} else if (s.Status.Water >= 6 && s.Status.Water <= 8) || (s.Status.Wind >= 7 && s.Status.Wind <= 15) {
		s.Status.DataStatus = "siaga"
	} else if s.Status.Water > 8 || s.Status.Wind > 15 {
		s.Status.DataStatus = "bahaya"
	}
	return &s
}

func GenerateToJson() {
	var datas entity.DataStatus
	for {
		datas.Status.Water = rand.Intn(100)
		datas.Status.Wind = rand.Intn(100)

		// write to json file
		jsonString, _ := json.Marshal(&datas)
		ioutil.WriteFile(jsonPath, jsonString, os.ModePerm)

		// sleep for 15 seconds
		time.Sleep(15 * time.Second)
	}
}
