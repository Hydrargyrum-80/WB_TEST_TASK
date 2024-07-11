package app

import (
	"WB_TEST_TASK/api"
	"WB_TEST_TASK/database"
	"WB_TEST_TASK/models"
	"encoding/json"
	"log"
	"sync"
	"time"
)

const APIKEY = "a40b0f15a41d1c6d3a99ddb82155fe94"

func UpdateCity(db database.Database) {
	log.Println("Cities update start!")
	citiesMap := make(map[string]int, 20)
	citiesMap["Moskow"] = 643
	citiesMap["Saint-Petersburg"] = 643
	citiesMap["Novosibirsk"] = 643
	citiesMap["Ekaterinburg"] = 643
	citiesMap["Kazan"] = 643
	citiesMap["Krasnoyarsk"] = 643
	citiesMap["Chelyabinsk"] = 643
	citiesMap["Ufa"] = 643
	citiesMap["Samara"] = 643
	citiesMap["Rostov-on-Don"] = 643
	citiesMap["Krasnodar"] = 643
	citiesMap["Omsk"] = 643
	citiesMap["Voronezh"] = 643
	citiesMap["Perm"] = 643
	citiesMap["Volgograd"] = 643
	citiesMap["Arkhangelsk"] = 643
	citiesMap["Astrakhan"] = 643
	citiesMap["Barnaul"] = 643
	citiesMap["Beloretsk"] = 643
	citiesMap["Blagoveshchensk"] = 643
	wg := sync.WaitGroup{}
	wg.Add(len(citiesMap))
	cityCh := make(chan models.CityModel)
	for key, value := range citiesMap {
		go func(key string, value int, out chan<- models.CityModel) {
			defer wg.Done()
			city, err := api.GetCityOpenWeatherAPI(key, value)
			if err != nil {
				log.Println(err)
				return
			}
			if city == nil {
				log.Println("City is nil!")
				return
			}
			out <- *city
		}(key, value, cityCh)
	}
	go func() {
		wg.Wait()
		close(cityCh)
	}()
	for city := range cityCh {
		err := db.AddCity(city)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("Cities update has been completed!")
}

func UpdateWeatherPredicts(db database.Database) {
	log.Println("Weather predict update start!")
	cityCh := make(chan models.CityModel)
	go func() {
		defer close(cityCh)
		err := db.GetAllCityChan(cityCh)
		if err != nil {
			log.Fatal(err)
		}
	}()
	predictCh := make(chan models.PredictInfoModel)
	wg := sync.WaitGroup{}
	for city := range cityCh {
		wg.Add(1)
		go func(city models.CityModel, output chan<- models.PredictInfoModel) {
			defer wg.Done()
			record, err := api.GetWeatherPredictOpenWeatherAPI(city)
			if err != nil {
				log.Println(err)
				return
			}
			for _, value := range record.List {
				el := value.(map[string]interface{})
				predict := models.PredictInfoModel{}
				predict.FkCityId = city.Id
				strDate := el["dt_txt"].(string)
				parseDate, err := time.Parse(time.DateTime, strDate)
				if err != nil {
					log.Println(err)
					return
				}
				predict.Date = parseDate
				mainMap := el["main"].(map[string]interface{})
				predict.Temp = mainMap["temp"].(float64)
				predict.Info, _ = json.Marshal(el)
				predictCh <- predict
			}
		}(city, predictCh)
	}
	go func() {
		wg.Wait()
		close(predictCh)
	}()
	var (
		updateCount uint
		err         error
	)
	for predict := range predictCh {
		updateCount, err = db.UpdateWeatherPredict(predict)
		if err != nil {
			log.Println(err)
			continue
		}
		if updateCount == 0 {
			err = db.AddWeatherPredict(predict)
			if err != nil {
				log.Println(err)
			}
		}
	}
	log.Println("Weather predict update has been completed!")
}
