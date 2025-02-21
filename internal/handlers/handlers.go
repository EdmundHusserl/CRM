package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EdmundHusserl/CRM/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Customer struct {
	Logger *logrus.Logger
	Repo   repository.CustomerRepository
}

type CustomerCreatedResponse struct {
	ID uuid.UUID `json:"id"`
}

type HandlerError struct {
	ErrorMsg string `json:"error_message"`
}

type CustomerHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

func NewCustomerHandler(logger *logrus.Logger, repo repository.CustomerRepository) CustomerHandler {
	return Customer{Logger: logger, Repo: repo}
}

// Create create a new customer
// @Summary Create a customer
// @Description Create customers
// @Accept  json
// @Produce  json
// @Param name query string true "Customer Name"
// @Param role query int true "Customer role" "Enum: 1=Base 2=Premium 3=Partner"
// @Param email query string true "Customer e-mail"
// @Param phone_number query string true "Customer phone number"
// @Param contacted query boolean true "Customer Contacted status"
// @Success 200 {object} CustomerCreatedResponse
// @Failure 400 {object} HandlerError
// @Failure 422 {object} HandlerError
// @Failure 500 {object} HandlerError
// @Router /api/customers [post]
func (h Customer) Create(w http.ResponseWriter, r *http.Request) {
	var c repository.Customer

	w.Header().Set("Content-Type", "application/json")
	jsonEnc := json.NewEncoder(w)

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Invalid e-mail format: %s", err.Error())}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"error_message": e.ErrorMsg,
			"status":        http.StatusBadRequest,
		}).Info("Failed to create new customer")
		return
	}
	c.ID = uuid.New()

	if _, err := c.ValidateEmail(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Invalid e-mail format: %s", c.Email)}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"error_message": e.ErrorMsg,
			"status":        http.StatusUnprocessableEntity,
		}).Info("Failed to create new customer")
		return
	}

	if err := c.ValidatePhone(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Invalid phone number format: %s", c.PhoneNumber)}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"error_message": e.ErrorMsg,
			"status":        http.StatusUnprocessableEntity,
		}).Info("Failed to create new customer")
		return
	}

	if err := h.Repo.Create(c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Could not create euser: %s", err.Error())}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"error_message": e.ErrorMsg,
			"status":        http.StatusInternalServerError,
		}).Warn("Failed to create new customer")
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonEnc.Encode(CustomerCreatedResponse{ID: c.ID})

	h.Logger.WithFields(logrus.Fields{
		"event":  fmt.Sprintf("ID: %v", c.ID),
		"status": http.StatusOK,
	}).Info("New record created")
}

// GetAll Get all customers
// @Summary Get all customers
// @Description Get all customers
// @Accept  json
// @Produce  json
// @Success 200 {object} []repository.Customer
// @Failure 500 {object} HandlerError
// @Router /api/customers [get]
func (h Customer) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonEnc := json.NewEncoder(w)

	customers, err := h.Repo.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Could not get users: %s", err.Error())}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"error_message": e.ErrorMsg,
			"status":        http.StatusInternalServerError,
		}).Warn("Failed to create new customer")

		return
	}
	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(customers)
}

// Get customer
// @Summary Get a customer by id
// @Description Get a customer by id
// @Accept  json
// @Produce  json
// @Param id query uuid.UUID true "User id"
// @Success 200 {object} repository.Customer
// @Failure 404 {object} HandlerError
// @Failure 422 {object} HandlerError
// @Router /api/customers/{id} [get]
func (h Customer) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonEnc := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Invalid user ID format: %s", id)}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"error_message": e.ErrorMsg,
			"status":        http.StatusUnprocessableEntity,
		}).Warn("Failed to create new customer")

		return
	}

	c, err := h.Repo.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		e := HandlerError{ErrorMsg: "User not found"}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"error_message": e.ErrorMsg,
			"status":        http.StatusNotFound,
		}).Warn("Failed to create new customer")

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(c)
}

// Delete customer
// @Summary Delete a customer
// @Description Get a customer
// @Accept  json
// @Produce  json
// @Param id query uuid.UUID true "User id"
// @Success 204 {object} nil
// @Failure 422 {object} HandlerError
// @Router /api/customers/{id} [delete]
func (h Customer) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonEnc := json.NewEncoder(w)

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Invalid user ID format: %s", id)}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"event":  fmt.Sprintf("ID: %v", id),
			"status": http.StatusUnprocessableEntity,
		}).Info("Deletion failure")

		return
	}

	err = h.Repo.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Could not delete user %s: %s", id, err.Error())}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"event":  fmt.Sprintf("ID: %v", id),
			"status": http.StatusBadRequest,
		}).Warn("Deletion failure")

		return
	}

	w.WriteHeader(http.StatusNoContent)

	h.Logger.WithFields(logrus.Fields{
		"event":  fmt.Sprintf("ID: %v", id),
		"status": http.StatusOK,
	}).Info("New record deleted")
}

// Update customer
// @Summary Update customer
// @Description Create customers
// @Accept  json
// @Produce  json
// @Param id query uuid.UUID true "User id"
// @Param name query string true "Customer Name"
// @Param role query int true "Customer role" "Enum: 1=Base 2=Premium 3=Partner"
// @Param email query string true "Customer e-mail"
// @Param phone_number query string true "Customer phone number"
// @Param contacted query boolean true "Customer Contacted status"
// @Success 200 {object} repository.Customer
// @Failure 400 {object} HandlerError
// @Failure 422 {object} HandlerError
// @Failure 500 {object} HandlerError
// @Router /api/customers [patch]
func (h Customer) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonEnc := json.NewEncoder(w)

	var c repository.Customer

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Invalid request payload: %s", err.Error())}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"event":  fmt.Sprintf("ID: %v", c.ID),
			"status": http.StatusBadRequest,
		}).Info("Update failure")

		return
	}

	if err := h.Repo.Update(c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e := HandlerError{ErrorMsg: fmt.Sprintf("Could not update user: %s", err.Error())}
		jsonEnc.Encode(e)

		h.Logger.WithFields(logrus.Fields{
			"event":  fmt.Sprintf("ID: %v", c.ID),
			"status": http.StatusBadRequest,
		}).Warn("Update failure")

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonEnc.Encode(c)
}
