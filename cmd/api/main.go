package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
)

func main() {

	loadEnv()

	config := config{
		address: env.GetEnvAsSting("PORT", "8080"),
		dbConfig: dbConfig{
			url:                env.GetEnvAsSting("DB_URL", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConnections: env.GetEnvAsInt("DB_MAX_OPEN_CONNECTIONS", 30),
			maxIdleConnections: env.GetEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 30),
			maxIdleTime:        env.GetEnvAsSting("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	conntection, err := db.CreateDbConnection(
		config.dbConfig.url,
		config.dbConfig.maxOpenConnections,
		config.dbConfig.maxIdleConnections,
		config.dbConfig.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer conntection.Close()

	log.Println("database connection pool established")

	postgresDb := db.PostgresDb(conntection)

	app := &application{
		config: config,
		db:     postgresDb,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
