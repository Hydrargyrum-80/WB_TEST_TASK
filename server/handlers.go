package server

import (
	"WB_TEST_TASK/models"
	"encoding/json"
	"net/http"
	"sort"
	"time"
)

func (s *Server) GetCityList(w http.ResponseWriter, r *http.Request) {
	var (
		response models.ListOfCityResponse = models.ListOfCityResponse{}
		err      error
	)
	response.Cities, err = s.storage.GetAllCityList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sort.Slice(response.Cities, func(i, j int) bool {
		return response.Cities[i].Name < response.Cities[j].Name
	})
	var buf []byte
	buf, err = json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetShortPredictCity(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("name")
	cityCountry := r.URL.Query().Get("country")
	city, err := s.storage.GetCity(cityName, cityCountry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var (
		predicts []models.PredictInfoModel
		response = models.ShortPredictResponse{CountryName: city.Country, CityName: city.Name}
	)
	predicts, err = s.storage.GetPredictsCity(*city, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var date time.Time
	for _, predict := range predicts {
		el := make(map[string]interface{})
		err = json.Unmarshal(predict.Info, &el)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mainMap := el["main"].(map[string]interface{})
		temp := mainMap["temp"].(float64)
		date, err = time.Parse(time.DateTime, el["dt_txt"].(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Elements = append(response.Elements, models.ShortPredictElem{Date: date, Temp: temp})
	}
	sort.Slice(response.Elements, func(i, j int) bool {
		return response.Elements[i].Date.Before(response.Elements[j].Date)
	})
	buf, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetFullPredictCityByTime(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("name")
	cityCountry := r.URL.Query().Get("country")
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	timeStr := r.URL.Query().Get("time")
	t, err := time.Parse(time.TimeOnly, timeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dateTime := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	city, err := s.storage.GetCity(cityName, cityCountry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var predict *models.PredictInfoModel
	predict, err = s.storage.GetPredictCity(*city, dateTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var infoMap map[string]interface{}
	err = json.Unmarshal(predict.Info, &infoMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := models.FullPredictResponse{CountryName: city.Country, CityName: city.Name, Temp: predict.Temp, Date: predict.Date, Info: infoMap}
	buf, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
