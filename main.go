package main

import (
	"log"

	"github.com/dunky-star/simply-bank/api"
	db "github.com/dunky-star/simply-bank/db/sqlc"
	"github.com/dunky-star/simply-bank/internal/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn := util.NewDB(config)
	store := db.NewStore(conn)

	server := api.NewServer(store)
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
