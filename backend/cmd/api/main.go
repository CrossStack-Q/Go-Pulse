package main

import (
	"log"

	"github.com/CrossStack-Q/Go-Pulse/backend/internal/db"
	"github.com/CrossStack-Q/Go-Pulse/backend/internal/env"
	"github.com/CrossStack-Q/Go-Pulse/backend/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8900"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:dsa@localhost:5433/titi?sslmode=disable"),
			maxopenConns: int64(env.GetInt("Max_Conn", 30)),
			maxIdleConn:  int64(env.GetInt("Max_Idle_Conn", 30)),
			maxIdleTime:  "3m",
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(
		cfg.db.addr,
		int(cfg.db.maxopenConns),
		int(cfg.db.maxopenConns),
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	log.Fatal(app.run(app.mount()))

}
