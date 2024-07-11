package models

import "time"

type ListOfCityResponse struct {
	Cities []CityModel `json:"list"`
}

type ShortPredictResponse struct {
	CountryName string             `json:"country"`
	CityName    string             `json:"city"`
	Elements    []ShortPredictElem `json:"predicts"`
}

type ShortPredictElem struct {
	Temp float64   `json:"temp"`
	Date time.Time `json:"dt_txt"`
}

type FullPredictResponse struct {
	CountryName string                 `json:"country"`
	CityName    string                 `json:"city"`
	Temp        float64                `json:"temp"`
	Date        time.Time              `json:"dt_txt"`
	Info        map[string]interface{} `json:"info"`
}
