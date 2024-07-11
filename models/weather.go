package models

import "time"

type CityModel struct {
	Id      int     `json:"-"`
	Name    string  `json:"name"`
	Lat     float32 `json:"lat"`
	Lon     float32 `json:"lon"`
	Country string  `json:"country"`
}

type PredictInfoModel struct {
	Id       int
	FkCityId int
	Temp     float64
	Date     time.Time
	Info     []byte
}

type PredictAPIResponseModel struct {
	List []interface{} `json:"list"`
}
