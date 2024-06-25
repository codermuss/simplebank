package main

import (
	"database/sql"
	"log"

	// * Note [codermuss]: Our code would not be able to talk to the databes without this import
	_ "github.com/lib/pq"
	"github.com/mustafayilmazdev/simplebank/api"
	db "github.com/mustafayilmazdev/simplebank/db/sqlc"
	"github.com/mustafayilmazdev/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)

	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
