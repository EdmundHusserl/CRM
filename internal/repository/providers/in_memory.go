package providers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/EdmundHusserl/CRM/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type InMemoryCustomerRepository struct {
	Customers []repository.Customer
}

func NewInMemoryCustomerRepository(data []repository.Customer) *InMemoryCustomerRepository {
	return &InMemoryCustomerRepository{Customers: data}
}

func (r *InMemoryCustomerRepository) CloseDBConnection() error {
	return nil
}

func (r *InMemoryCustomerRepository) Create(c repository.Customer) error {
	for _, customer := range r.Customers {
		if c.ID == customer.ID {
			return fmt.Errorf("conflict: user %s does exist", c.ID)
		}
		if c.Email == customer.Email {
			return fmt.Errorf("conflict: email %s does exist", c.Email)
		}
	}

	r.Customers = append(r.Customers, c)
	return nil
}

func (r *InMemoryCustomerRepository) Get(id uuid.UUID) (*repository.Customer, error) {
	for _, c := range r.Customers {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("user not found: %v", id)
}

func (r *InMemoryCustomerRepository) Delete(id uuid.UUID) error {
	var index int
	found := false
	for i, c := range r.Customers {
		if c.ID == id {
			found = true
			index = i
			break
		}
	}
	if !found {
		return fmt.Errorf("user not found: %v", id)
	}
	customers := make([]repository.Customer, len(r.Customers)-1)
	customers = append(r.Customers[:index], r.Customers[index+1:]...)
	r.Customers = customers
	return nil
}

func (r *InMemoryCustomerRepository) GetAll() ([]repository.Customer, error) {
	return r.Customers, nil
}

func (r *InMemoryCustomerRepository) Update(c repository.Customer) error {
	for i, customer := range r.Customers {
		if customer.ID == c.ID {
			r.Customers[i].Name = c.Name
			r.Customers[i].Role = c.Role
			r.Customers[i].Email = c.Email
			r.Customers[i].PhoneNumber = c.PhoneNumber
			r.Customers[i].Contacted = c.Contacted
			return nil
		}
	}
	return fmt.Errorf("could not update user with ID %v", c.ID)
}

func LoadFromCSVFile(l *logrus.Logger, path string) ([]repository.Customer, error) {
	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		l.WithField(
			"event",
			fmt.Sprintf("opening csv file name %s threw error: %s", path, err.Error()),
		).Error("failure while opening files")
		return nil, err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		l.WithField("event", err.Error()).Error("failure during data parsing")
		return nil, err
	}

	customers := []repository.Customer{}

	// Loop through lines (skipping header)
	for i, line := range lines {
		if i == 0 {
			continue // Skip header row
		}

		role, err := strconv.Atoi(line[1])
		if err != nil {
			l.WithField("event", err.Error()).Warn(fmt.Sprintf("failed parsing role column %s", line[1]))
		}
		contacted, err := strconv.ParseBool(line[4])
		if err != nil {
			l.WithField("event", err.Error()).Warn(fmt.Sprintf("failed parsing contacted column %s", line[4]))
		}

		c := repository.Customer{
			ID:          uuid.New(),
			Name:        line[0],
			Role:        repository.CustomerRole(role),
			Email:       line[2],
			PhoneNumber: line[3],
			Contacted:   contacted,
		}
		customers = append(customers, c)
	}

	return customers, nil
}
