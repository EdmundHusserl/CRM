package providers

import (
	"fmt"
	"strings"

	"github.com/EdmundHusserl/CRM/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	psql      string = "psql"
	in_memory string = "in-memory"
)

func isValid(provider string) bool {
	return (provider == psql || provider == in_memory)
}

// Returns a CustomerRepository implemented interface
func NewRepository(l *logrus.Logger, provider string) repository.CustomerRepository {
	if !isValid(provider) {
		l.WithField(
			"event", fmt.Sprintf("defaulting to %s", in_memory),
		).Info("Unknown DB provider")
	}
	switch strings.ToLower(provider) {
	case "psql":
		return NewPostgresCustomerRepository(l)
	default:
		return &InMemoryCustomerRepository{Customers: []repository.Customer{}}
	}
}
