package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/EdmundHusserl/CRM/internal/handlers"
	psql "github.com/EdmundHusserl/CRM/internal/repository/providers"
	"github.com/EdmundHusserl/CRM/internal/router"

	_ "github.com/lib/pq"
)

func main() {

	connStr := psql.GetConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening psql instance db: %v", err)
	}
	defer db.Close()
	repo := psql.NewPostgresCustomerRepository(db)
	handler := handlers.NewCustomerHandler(repo)	
	router := router.NewRouter(handler) 
  fmt.Println("Server is starting on port 3000...")
  http.ListenAndServe(":3000", router)
}