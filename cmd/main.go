package main

import (
	"flag"

	"github.com/EdmundHusserl/CRM/internal/server"
)

func main() {
	dbProvider := flag.String("db", "in-memory", "DB provider: in-memory|psql")
	serverPort := flag.Int("port", 3000, "Server port")
	flag.Parse()

	server := server.NewServer(*dbProvider, *serverPort)
	defer server.DB.CloseDBConnection()

	server.Listen()
}
