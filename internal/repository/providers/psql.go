package providers

import (
	"database/sql"
	"errors"

	"github.com/EdmundHusserl/CRM/internal/repository"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresCustomerRepository struct {
    db *sql.DB
}

func NewPostgresCustomerRepository(db *sql.DB) *PostgresCustomerRepository {
  return &PostgresCustomerRepository{db: db}
}

func (r *PostgresCustomerRepository) Create(c repository.Customer) error {

	if _, err := c.ValidateEmail(); err != nil {
		return err
	} 
	
	if err := c.ValidatePhone(); err != nil {
		return err
	}
	_, err := r.db.Exec(
			"INSERT INTO customers (id, name, role, email, phone_number) VALUES ($1, $2, $3, $4, $5)", 
			c.ID, c.Name, c.Email, c.Role, c.PhoneNumber)
  return err
}

func (r *PostgresCustomerRepository) Get(id uuid.UUID) (*repository.Customer, error) {
  c := &repository.Customer{}
  err := r.db.QueryRow(
		"SELECT id, name, role, email, phone_number FROM customers WHERE id=$1", id).Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.PhoneNumber)
  
		if err == sql.ErrNoRows {
      return nil, errors.New("user not found")
  }
  return c, err
}

func (r *PostgresCustomerRepository) GetAll() ([]repository.Customer, error) {
  rows, err := r.db.Query("SELECT id, name, role, email, phone_number FROM customers")
  if err != nil {
      return nil, err
  }
  defer rows.Close()
  var customers []repository.Customer
  for rows.Next() {
      var c repository.Customer
      if err := rows.Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.PhoneNumber); err != nil {
          return nil, err
      }
      customers = append(customers, c)
  }
  return customers, nil
}

func (r *PostgresCustomerRepository) Delete(id uuid.UUID) error {
  _, err := r.db.Exec("DELETE FROM customers WHERE id=$1", id)
  return err
}
