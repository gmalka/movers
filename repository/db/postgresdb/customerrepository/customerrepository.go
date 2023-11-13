package customerrepository

import (
	"fmt"

	"github.com/gmalka/movers/model"
	"github.com/jmoiron/sqlx"
)

type CustomerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (c *CustomerRepository) CreateCustomer(customer model.CustomerInfo) error {
	_, err := c.db.Exec("INSERT INTO customers VALUES($1,$2)", customer.Name, customer.Money)
	if err != nil {
		return fmt.Errorf("cant insert customer %s: %v", customer.Name, err)
	}

	return nil
}

func (c *CustomerRepository) GetCustomer(name string) (model.CustomerInfo, error) {
	row := c.db.QueryRow("SELECT * FROM customers WHERE name = $1", name)

	customer := model.CustomerInfo{}
	err := row.Scan(&customer.Name, &customer.Money)
	if err != nil {
		return model.CustomerInfo{}, fmt.Errorf("cant find customer %s: %v", name, err)
	}

	return customer, err
}