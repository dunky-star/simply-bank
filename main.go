package main

import (
	"log"

	"github.com/dunky-star/simply-bank/api"
	db "github.com/dunky-star/simply-bank/db/sqlc"
	"github.com/dunky-star/simply-bank/internal/util"
)

const serverAddress = "0.0.0.0:8080"

func main() {
	conn := util.NewDB()
	store := db.NewStore(conn)

	server := api.NewServer(store)
	if err := server.Start(serverAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
