package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {

	loadEnv()

	config := config{
		address: env.GetEnvAsSting("PORT", "8080"),
		dbConfig: dbConfig{
			url:                env.GetEnvAsSting("DB_URL", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConnections: env.GetEnvAsInt("DB_MAX_OPEN_CONNECTIONS", 30),
			maxIdleConnections: env.GetEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 30),
			maxIdleTime:        env.GetEnvAsSting("DB_MAX_IDLE_TIME", "15m"),
			apiUrl:             env.GetEnvAsSting("EXTERNAL_URL", "localhost:8080"),
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
