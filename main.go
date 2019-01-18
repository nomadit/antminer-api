package main

import (
	"context"
	"flag"
	"github.com/nomadit/antminer-api/api"
	"github.com/nomadit/antminer-api/api/models/db"
	"github.com/nomadit/antminer-api/config"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	// ServiceMode is the mode of service
	ServiceMode = flag.String("mode", "dev", "The service mode")

	// ServicePort is the port
	ServicePort = flag.String("port", ":3100", "The service port")
	ConfigFilePath = flag.String("config", "antminer-api.yaml", "the service config")
)

func main() {
	flag.Parse()

	confMap := getConf()

	conf := (*confMap)[*ServiceMode]
	initDB(&conf.DB)

	runServer()
}

func getConf() *map[string]config.Config {
	viper.SetConfigFile(*ConfigFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	conf := map[string]config.Config{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return &conf
}

func initDB(conf *config.DBConfig)  {
	db := db.NewDB(conf)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(100)
}

func runServer()  {
	srv := &http.Server{
		Addr:    *ServicePort,
		Handler: api.RouteHandler(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}