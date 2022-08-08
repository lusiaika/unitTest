package handler

import (
	"autoreloaddata/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data entity.DataStatus

const jsonpath = "html/status.json"
const htmlpath = "html/web.html"

func Test_GenerateToJson(t *testing.T) {
	var data entity.DataStatus
	t.Log("Test write json file")
	go GenerateToJson()
	os.Chdir("..")
	cwd, _ := os.Getwd()
	fmt.Println(cwd)
	assert.FileExists(t, jsonpath)
	file, err := ioutil.ReadFile(jsonpath)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, file)

	err = json.Unmarshal(file, &data)
	assert.NotNil(t, data)
	if err != nil {
		assert.Error(t, err)
	}

}

func Test_checkstatus(t *testing.T) {
	t.Run("test_checkstatus", func(t *testing.T) {
		var data entity.DataStatus

		for i := 0; i <= 5; i++ {
			for j := 0; j <= 6; j++ {
				data.Status.Water = i
				data.Status.Wind = j
				got := CheckStatus(data)
				if got.Status.DataStatus != "Aman" {
					t.Error("Wrong Status")
				}
			}

		}

		for i := 6; i <= 8; i++ {
			for j := 7; j <= 15; j++ {
				data.Status.Water = i
				data.Status.Wind = j
				got := CheckStatus(data)
				if got.Status.DataStatus != "siaga" {
					t.Error("Wrong Status")
				}
			}
		}

		for i := 9; i <= 100; i++ {
			for j := 16; j <= 100; j++ {
				data.Status.Water = i
				data.Status.Wind = j
				got := CheckStatus(data)
				if got.Status.DataStatus != "bahaya" {
					t.Error("Wrong Status")
				}
			}
		}
	})
}
