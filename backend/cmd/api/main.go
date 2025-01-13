package main

import (
	"log"

	"github.com/CrossStack-Q/Go-Pulse/backend/internal/env"
	"github.com/CrossStack-Q/Go-Pulse/backend/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8900"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	log.Fatal(app.run(app.mount()))

}
