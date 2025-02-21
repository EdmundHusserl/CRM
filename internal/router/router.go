package router

import (
	"net/http"

	_ "github.com/EdmundHusserl/CRM/docs"
	"github.com/EdmundHusserl/CRM/internal/handlers"
	"github.com/gorilla/mux"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h handlers.CustomerHandler) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/api/customers/{id}", h.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/api/customers/{id}", h.Get).Methods(http.MethodGet)
	router.HandleFunc("/api/customers", h.Update).Methods(http.MethodPatch)
	router.HandleFunc("/api/customers", h.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/customers", h.GetAll).Methods(http.MethodGet)
	return router
}
