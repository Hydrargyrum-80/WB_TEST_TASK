package api

import (
	"WB_TEST_TASK/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const APIKEY = "a40b0f15a41d1c6d3a99ddb82155fe94"

func GetCityOpenWeatherAPI(cityName string, cityCode int) (*models.CityModel, error) {
	var (
		err  error
		resp *http.Response
		body []byte
	)
	resp, err = http.Get("http://api.openweathermap.org/geo/1.0/direct?q=" + cityName + "," + strconv.Itoa(cityCode) + "&limit=1&appid=" + APIKEY)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Println(resp.Status, " ", resp)
		return nil, nil
	}
	body, _ = io.ReadAll(resp.Body)
	var record []models.CityModel
	err = json.Unmarshal(body, &record)
	if err != nil {
		return nil, err
	}
	return &record[0], nil
}

func GetWeatherPredictOpenWeatherAPI(city models.CityModel) (*models.PredictAPIResponseModel, error) {
	var (
		err  error
		resp *http.Response
		body []byte
	)
	resp, err = http.Get("http://api.openweathermap.org/data/2.5/forecast?lat=" + fmt.Sprintf("%f", city.Lat) + "&lon=" + fmt.Sprintf("%f", city.Lon) + "&appid=" + APIKEY + "&units=metric")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		return nil, nil
	}
	body, _ = io.ReadAll(resp.Body)
	record := models.PredictAPIResponseModel{}
	record.List = make([]interface{}, 40)
	err = json.Unmarshal(body, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
