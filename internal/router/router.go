package router

import (
	"net/http"

	"github.com/EdmundHusserl/CRM/internal/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(h handlers.CustomerHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/customers/{id}", h.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/customers/{id}", h.Get).Methods(http.MethodGet)
	router.HandleFunc("/customers", h.Update).Methods(http.MethodPatch)
	router.HandleFunc("/customers", h.Create).Methods(http.MethodPost)
	router.HandleFunc("/customers", h.GetAll).Methods(http.MethodGet)
	return router 
}