package main

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "postgres"
)

func main() {
	params := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	server := NewApiServer(
		NewGameInfoService(),
		NewDetermineInfoService(),
		NewDataBaseService(sql.Open("postgres", params)),
	)
	server.Start()
}
