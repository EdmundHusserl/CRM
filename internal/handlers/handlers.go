package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EdmundHusserl/CRM/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Customer struct {
  Repo repository.CustomerRepository
}

type CustomerCreatedResponse struct {
	ID uuid.UUID `json:"id"`
}

type CustomerHandler interface {
	Create(w http.ResponseWriter, r *http.Request)  
	Delete(w http.ResponseWriter, r *http.Request) 
	Get(w http.ResponseWriter, r *http.Request) 
	GetAll(w http.ResponseWriter, r *http.Request)
}

func NewCustomerHandler(repo repository.CustomerRepository) CustomerHandler {
  return Customer{Repo: repo}
}

func (h Customer) Create(w http.ResponseWriter, r *http.Request) {
  var c repository.Customer
  w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
      http.Error(w, "Invalid request payload", http.StatusBadRequest)
      return
  }
  c.ID = uuid.New()
  if err := h.Repo.Create(c); err != nil {
      http.Error(w, "Could not create user", http.StatusInternalServerError)
      return
  }
  w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CustomerCreatedResponse{ID: c.ID})
}

func (h Customer) GetAll(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    customers, err := h.Repo.GetAll()
		if err != nil {
      http.Error(w, "Could not get users", http.StatusInternalServerError)
      return
    }
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customers)
}


func (h Customer) Get(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		id, err := uuid.Parse(vars["id"])  
		if err != nil {	
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)	
		}
		c, err := h.Repo.Get(id) 
		if err != nil {
			http.Error(w, "Could not get user", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(c)
}

func (h Customer) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])  
	if err != nil {	
		http.Error(w, "Invalid id format", http.StatusBadRequest)	
	}
	err = h.Repo.Delete(id)
	if err != nil {
		http.Error(w, "Could not get user", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

