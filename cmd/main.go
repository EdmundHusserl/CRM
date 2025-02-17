package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/EdmundHusserl/CRM/internal/handlers"
	"github.com/EdmundHusserl/CRM/internal/repository/providers"
	"github.com/EdmundHusserl/CRM/internal/router"

	_ "github.com/lib/pq"
)

func main() {
   
	connStr := "host=localhost port=5432 user=postgres password=p4ssw0rd dbname=customers sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening psql instance db: %v", err)
	}
	defer db.Close()
	repo := providers.NewPostgresCustomerRepository(db)
	handler := handlers.NewCustomerHandler(repo)	
	router := router.NewRouter(handler) 
  fmt.Println("Server is starting on port 3000...")
  http.ListenAndServe(":3000", router)
}