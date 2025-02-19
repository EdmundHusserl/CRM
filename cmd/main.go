package main

import (
	"flag"
	"fmt"

	"github.com/EdmundHusserl/CRM/internal/server"
)

func main() {
	dbProvider := flag.String("db", "in-memory", "DB provider: in-memory|psql. Default: in-memory.")
	serverPort := flag.Int("port", 3000, "Server port. Default: 3000")
	flag.Parse()

	server := server.NewServer(*dbProvider, *serverPort)
	defer server.DB.CloseDBConnection()

	server.Logger.WithField(
		"event",
		fmt.Sprintf("Server is starting on port %s...", server.Addr[1:]),
	)
	server.Listen()
}
