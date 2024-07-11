package database

import (
	"WB_TEST_TASK/models"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"strings"
	"time"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	db, err := sql.Open("postgres", "user=postgres password=2917819 dbname=WeatherDB sslmode=disable")
	if err != nil {
		panic(err)
	}
	return &Database{db: db}
}

func (db *Database) CreateTables() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	migrateFilePath := dir + "/database/migrate.sql"
	migrateFile, err := os.ReadFile(migrateFilePath)
	if err != nil {
		return err
	}
	createTableRequests := strings.Split(string(migrateFile), ";")
	for _, request := range createTableRequests {
		_, err = db.db.Exec(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) AddCity(newCity models.CityModel) error {
	query := "INSERT INTO city (city_name, city_lat, city_lon, city_country) VALUES ($1, $2, $3, $4)"
	_, err := db.db.Exec(query, newCity.Name, newCity.Lat, newCity.Lon, newCity.Country)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetAllCityChan(outputChan chan<- models.CityModel) error {
	rows, err := db.db.Query("SELECT city_id, city_name, city_lat, city_lon, city_country FROM city")
	if err != nil {
		return err
	}
	defer rows.Close()
	var city models.CityModel
	for rows.Next() {
		city = models.CityModel{}
		err = rows.Scan(&city.Id, &city.Name, &city.Lat, &city.Lon, &city.Country)
		if err != nil {
			return err
		}
		outputChan <- city
	}
	return nil
}

func (db *Database) GetAllCityList() ([]models.CityModel, error) {
	rows, err := db.db.Query("SELECT city_id, city_name, city_lat, city_lon, city_country FROM city")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		cities []models.CityModel
		city   models.CityModel
	)
	for rows.Next() {
		err = rows.Scan(&city.Id, &city.Name, &city.Lat, &city.Lon, &city.Country)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	return cities, nil
}

func (db *Database) GetCity(cityName string, countryName string) (*models.CityModel, error) {
	rows, err := db.db.Query("SELECT city_id, city_name, city_lat, city_lon, city_country FROM city WHERE city_name = $1 AND city_country = $2", cityName, countryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		city  models.CityModel
		count uint = 0
	)
	for rows.Next() {
		if count > 0 {
			panic("more than one city found")
		}
		err = rows.Scan(&city.Id, &city.Name, &city.Lat, &city.Lon, &city.Country)
		if err != nil {
			return nil, err
		}
		count += 1
	}
	return &city, nil
}

func (db *Database) GetPredictsCity(city models.CityModel, startDate time.Time) ([]models.PredictInfoModel, error) {
	rows, err := db.db.Query("SELECT temp, predict_date, info FROM weather_predict WHERE FK_city_id = $1 AND predict_date > $2", city.Id, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		predict models.PredictInfoModel
		list    []models.PredictInfoModel
	)
	for rows.Next() {
		predict = models.PredictInfoModel{}
		err = rows.Scan(&predict.Temp, &predict.Date, &predict.Info)
		if err != nil {
			return nil, err
		}
		list = append(list, predict)
	}
	return list, nil
}

func (db *Database) GetPredictCity(city models.CityModel, Date time.Time) (*models.PredictInfoModel, error) {
	rows, err := db.db.Query("SELECT temp, predict_date, info FROM weather_predict WHERE FK_city_id = $1 AND predict_date = $2", city.Id, Date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		predict models.PredictInfoModel
		count   uint = 0
	)
	for rows.Next() {
		if count > 1 {
			panic("more than one predict found")
		}
		err = rows.Scan(&predict.Temp, &predict.Date, &predict.Info)
		if err != nil {
			return nil, err
		}
		count++
	}
	return &predict, nil
}

func (db *Database) UpdateWeatherPredict(predict models.PredictInfoModel) (uint, error) {
	var (
		query        string = "UPDATE weather_predict SET temp = $1, info = $2 WHERE fk_city_id = $3 AND predict_date = $4"
		updateResult sql.Result
		err          error
	)
	updateResult, err = db.db.Exec(query, predict.Temp, predict.Info, predict.FkCityId, predict.Date)
	if err != nil {
		return 0, err
	}
	updateRowsCount, _ := updateResult.RowsAffected()
	return uint(updateRowsCount), nil
}

func (db *Database) AddWeatherPredict(predict models.PredictInfoModel) error {
	var (
		query string = "INSERT INTO weather_predict (fk_city_id, temp, predict_date, info) VALUES ($1, $2, $3, $4)"
		err   error
	)
	_, err = db.db.Exec(query, predict.FkCityId, predict.Temp, predict.Date, predict.Info)
	if err != nil {
		return err
	}
	return nil
}
