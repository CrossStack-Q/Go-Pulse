package main

import (
	"log"

	"github.com/CrossStack-Q/Go-Pulse/backend/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8900"),
	}
	app := &application{
		config: cfg,
	}

	log.Fatal(app.run(app.mount()))

}
