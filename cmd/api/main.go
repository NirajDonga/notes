package main

import (
	"crud/internal/config"
	"crud/internal/db"
	"crud/internal/server"
	"fmt"
	"log"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config Error")
	}

	client, database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db Error %v", err)
	}

	defer func() {
		if err := db.DisConnect(client); err != nil {
			log.Printf("Mongo Disconnected error: %v", err)
		}
	}()

	router := server.NewRouter(database)

	addr := fmt.Sprintf(":%s", cfg.ServerPORT)
	if err := router.Run(addr); err != nil {
		log.Fatal("Server failed")
	}

}
