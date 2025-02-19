package providers

import (
	"testing"

	"github.com/EdmundHusserl/CRM/internal/repository"
	"github.com/google/uuid"
)


func TestCreateCustomer(t *testing.T) {
	cUUID := uuid.New()
		
	tests := []struct{
		name 		string
		repo    *InMemoryCustomerRepository
		data    repository.Customer
		count   int
		wantErr bool
	}{
		
  	{
			"Success", 
			&InMemoryCustomerRepository{Customers: []repository.Customer{}}, 
			repository.Customer{
				ID: cUUID,
				Name: "Jorge",
				Role: 2,
				Email: "jorge@corp.com",
				PhoneNumber: "514 888 8888",
			},
			1, 
			false,
		},
    {
			"Fails due to ID conflict", 
			&InMemoryCustomerRepository{Customers: []repository.Customer{
				{
					ID: cUUID,
					Name: "Jorge",
					Role: 1,
					Email: "jorge@corp.com",
					PhoneNumber: "514 888 8888",	
				},
			}}, 
			repository.Customer{
				ID: cUUID,
				Name: "Whatever Dude",
				Role: 1,
				Email: "whatdud@corp.com",
				PhoneNumber: "514 999 8888",
			},
			1, 
			true,
		},
		{
			"Fails due to e-mail conflict", 
			&InMemoryCustomerRepository{Customers: []repository.Customer{
				{
					ID: cUUID,
					Name: "Jorge",
					Role: 1,
					Email: "jorge@corp.com",
					PhoneNumber: "514 888 8888",	
				},
			}}, 
			repository.Customer{
				ID: uuid.New(),
				Name: "Whatever Dude",
				Role: 2,
				Email: "jorge@corp.com",
				PhoneNumber: "514 999 7777",
			},
			1, 
			true,
		},
		}
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      
			err := tt.repo.Create(tt.data)
      if (err != nil) != tt.wantErr {
          t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
      }
			
			customerCount := len(tt.repo.Customers)
			if customerCount != tt.count {
				t.Errorf("Customer count=%v not equal to expected=%v", customerCount, tt.count)
			}
    
		})
  }	
}

func TestGetDeleteCustomer(t *testing.T) {
	cUUID := uuid.New()
		
	tests := []struct{
		name 		string
		repo 		*InMemoryCustomerRepository
		data    uuid.UUID
		count   map[string]int
		wantErr bool
	}{
		
  	{
			"Success", 
			&InMemoryCustomerRepository{Customers: []repository.Customer{
				{
					ID: cUUID,
					Name: "Jorge",
					Role: 2,
					Email: "jorge@corp.com",
					PhoneNumber: "514 888 8888",
				},
			}},
			cUUID, 
			map[string]int{"GET": 1, "DELETE": 0},
			false,
		},
    {
			"Fails_as_not_found", 
			&InMemoryCustomerRepository{Customers: []repository.Customer{
				{
					ID: cUUID,
					Name: "Jorge",
					Role: 1,
					Email: "jorge@corp.com",
					PhoneNumber: "514 888 8888",	
				},
			}},
			uuid.New(),
			map[string]int{"GET": 1, "DELETE": 1},
			true,
		},
	}
	
	for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      
			_, err := tt.repo.Get(tt.data)
      
			if (err != nil) != tt.wantErr {
        t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
      }
			
			customerCount := len(tt.repo.Customers)
			countAfterGet, _ := tt.count["GET"]
			if customerCount != countAfterGet {
				t.Errorf("Customer count=%v not equal to expected=%v", customerCount, countAfterGet)
			}
			
			err = tt.repo.Delete(tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error =%v, wantErr %v", err, tt.wantErr)
			}

			customerCount = len(tt.repo.Customers)
			countAfterDelete, _ := tt.count["DELETE"] 
			if customerCount != countAfterDelete {
				t.Errorf("Customer count=%v not equal to expected=%v", customerCount, countAfterDelete)
			}

		})
  }
		
}

func TestGetAllCustomers(t *testing.T) {
	cUUID := uuid.New()
	
	tests := []struct{
		name 		string
		repo 		*InMemoryCustomerRepository
		count   int
		wantErr bool
	}{
		
  	{
			"Success", 
			&InMemoryCustomerRepository{Customers: []repository.Customer{
				{
					ID: cUUID,
					Name: "Jorge",
					Role: 2,
					Email: "jorge@corp.com",
					PhoneNumber: "514 888 8888",
				},
			}}, 
			1,
			false,
		},
	}
	
	for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      
			customers, err := tt.repo.GetAll()
      
			if (err != nil) != tt.wantErr {
        t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
      }
			
			customerCount := len(tt.repo.Customers)
			countAfterRetrieval := len(customers)
			if customerCount != countAfterRetrieval {
				t.Errorf("Customer count=%v not equal to expected=%v", customerCount, countAfterRetrieval)
			}
			
		})
  }
}

func TestUpdateCustomer(t *testing.T) {
	cUUID := uuid.New()
	
	tests := []struct{
		name 		string
		repo 		*InMemoryCustomerRepository
		data    repository.Customer
		wantErr bool
	}{
		
  	{
			"Success", 
			&InMemoryCustomerRepository{Customers: []repository.Customer{
				{
					ID: cUUID,
					Name: "Jorge",
					Role: 2,
					Email: "jorge@corp.com",
					PhoneNumber: "514 888 8888",
				},
			}},
			repository.Customer{
				ID: cUUID,
				Name: "Jorge",
				Role: 2,
				Email: "jorge@corp.com",
				PhoneNumber: "514 888 8888",
			},
			false,
		},
		{
			"Fails_as_not_found", 
			&InMemoryCustomerRepository{Customers: []repository.Customer{
				{
					ID: cUUID,
					Name: "Jorge",
					Role: 2,
					Email: "jorge@corp.com",
					PhoneNumber: "514 888 8888",
				},
			}},
			repository.Customer{
				ID: uuid.New(),
				Name: "Jorge",
				Role: 1,
				Email: "jorge@corp.com",
				PhoneNumber: "514 777 8888",
			},
			true,
		},
	}
	
	for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      
			err := tt.repo.Update(tt.data)
			if (err != nil) != tt.wantErr {
        t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
      }
			
		})
  }	
}