package main

import (
	"database/sql"
	"log"

	// * Note [codermuss]: Our code would not be able to talk to the databes without this import
	_ "github.com/lib/pq"
	"github.com/mustafayilmazdev/simplebank/api"
	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
