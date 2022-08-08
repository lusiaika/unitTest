package entity

type DataStatus struct {
	Status Data `json:"status"`
}

type Data struct {
	Water      int    `json:"water"`
	Wind       int    `json:"wind"`
	DataStatus string `json:"datastatus,omitempty"`
}
