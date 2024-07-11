package server

import (
	"WB_TEST_TASK/database"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	storage *database.Database
}

func NewServer(storage *database.Database) *Server {
	return &Server{storage: storage}
}

func (s *Server) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	r := router.With(s.MeasureRequestTime)
	r.Get("/get_city_list", s.GetCityList)
	r.Get("/full_predict_city", s.GetFullPredictCityByTime)
	r.Get("/short_predict_city", s.GetShortPredictCity)
	return router
}
