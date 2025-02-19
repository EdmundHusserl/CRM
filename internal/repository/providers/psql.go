package providers

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/EdmundHusserl/CRM/internal/repository"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type PostgresCustomerRepository struct {
	db *sql.DB
}

func getEnvOrDefault(envVarName, defaultTo string) string {
	envVar := os.Getenv(envVarName)
	if len(envVar) == 0 {
		return defaultTo
	}
	return envVar
}

// Produces connection string
func getConnectionString() string {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	password := os.Getenv("DB_PASSWORD")
	dbName := getEnvOrDefault("DB_NAME", "customers")
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, port, dbUser, password, dbName)

	return connStr
}

func NewPostgresCustomerRepository(l *logrus.Logger) *PostgresCustomerRepository {
	connStr := getConnectionString()

	l.WithField("event", fmt.Sprintf("attempting psql connection with %s", connStr)).Info("db connection")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.WithField(
			"error",
			err.Error(),
		).Fatal("error opening psql instance")
	}

	return &PostgresCustomerRepository{db: db}
}

func (r *PostgresCustomerRepository) CloseDBConnection() error {
	return r.db.Close()
}

func (r *PostgresCustomerRepository) Create(c repository.Customer) error {
	_, err := r.db.Exec(
		"INSERT INTO customers (id, name, role, email, phone_number) VALUES ($1, $2, $3, $4, $5)",
		c.ID, c.Name, c.Role, c.Email, c.PhoneNumber)
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

func (r *PostgresCustomerRepository) Update(c repository.Customer) error {
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Printf("DB operational error: %v\n", err)
		return err
	}
	query := "UPDATE customers SET name=$2, role=$3, email=$4, phone_number=$5 WHERE id=$1"
	_, err = tx.Exec(query, c.ID, c.Name, c.Role, c.Email, c.PhoneNumber)
	if err != nil {
		fmt.Printf("DB operational error: %v\n", err)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		fmt.Printf("DB operational error: %v\n", err)
		tx.Rollback()
		return err
	}
	fmt.Println("Transaction commited successfully")
	return nil
}
