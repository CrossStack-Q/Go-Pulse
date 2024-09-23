package main

import (
	"log"

	"github.com/CrossStack-Q/Go-Pulse/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":9090"),
	}

	app := &application{
		config: cfg,
	}

	// env load
	//	os.LookupEnv("PATH")

	mux := app.mount()

	log.Fatal(app.run(mux))
}
