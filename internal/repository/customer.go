package repository

import (
	"fmt"
	"net/mail"
	"regexp"

	"github.com/google/uuid"
)

// Define a new Enum which is descriptive of client roles.
type customerRole int

const (
	Basic = iota
	Premium
	Partner
)

type Customer struct {
	ID 					uuid.UUID
	Name 				string
	Role 			  customerRole
	Email 		  string 
	PhoneNumber string
}

type CustomerRepository interface {
	Create(c Customer) error
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*Customer, error)
	GetAll() ([]Customer, error) 
}

func (c Customer) ValidateEmail() (*string, error) {
	a, err := mail.ParseAddress(c.Email)

	if err != nil {
		return nil, err
	}

	return &a.Address, nil
}

func (c Customer) ValidatePhone() error {
	phoneRegex := regexp.MustCompile(`^\+?[0-9\s\-\(\)]+$`)
	
	if ok := phoneRegex.MatchString(c.PhoneNumber); !ok {
		return fmt.Errorf("invalid phone number format: %s", c.PhoneNumber)
	}
	
	return nil
}