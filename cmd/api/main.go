package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
)

func main() {

	loadEnv()

	postgresDb := db.PostgresDb(nil)

	config := config{
		address: env.GetEnvAsSting("PORT", "8080"),
	}

	app := &application{
		config: config,
		db:     postgresDb,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
