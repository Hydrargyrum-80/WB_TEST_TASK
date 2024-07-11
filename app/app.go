package app

import (
	"WB_TEST_TASK/database"
	"WB_TEST_TASK/server"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const MYTIME = 1

func Start() {
	db := database.NewDatabase()
	err := db.CreateTables()
	if err != nil {
		log.Fatal(err)
	}
	UpdateCity(*db)
	UpdateWeatherPredicts(*db)
	webserver := server.NewServer(db)
	router := webserver.InitRouter()
	go func() {
		for {
			time.Sleep(MYTIME * time.Minute)
			UpdateWeatherPredicts(*db)
		}
	}()
	httpServer := http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	<-signalChan
	shutdown := context.Background()
	if err = httpServer.Shutdown(shutdown); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("App has shut down")
}
