package customerrepository

import (
	"context"
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

func (c *CustomerRepository) CreateCustomer(ctx context.Context, customer model.CustomerInfo) error {
	_, err := c.db.ExecContext(ctx, "INSERT INTO customers VALUES($1,$2)", customer.Name, customer.Money)
	if err != nil {
		return fmt.Errorf("cant insert customer %s: %v", customer.Name, err)
	}

	return nil
}

func (c *CustomerRepository) UpdateCustomer(ctx context.Context, customer model.CustomerInfo) error {
	_, err := c.db.ExecContext(ctx, "UPDATE customers SET name = $1, money = $2", customer.Name, customer.Money)
	if err != nil {
		return fmt.Errorf("cant update customer %s: %v", customer.Name, err)
	}

	return nil
}

func (c *CustomerRepository) DeleteCustomer(ctx context.Context, name string)  error {
	_, err := c.db.ExecContext(ctx, "DELETE FROM customers WHERE name=$1", name)
	if err != nil {
		return fmt.Errorf("cant delete customer %s: %v", name, err)
	}

	return nil
}

func (c *CustomerRepository) GetCustomer(ctx context.Context, name string) (model.CustomerInfo, error) {
	row := c.db.QueryRowContext(ctx, "SELECT * FROM customers WHERE name = $1", name)

	customer := model.CustomerInfo{}
	err := row.Scan(&customer.Name, &customer.Money)
	if err != nil {
		return model.CustomerInfo{}, fmt.Errorf("cant find customer %s: %v", name, err)
	}

	return customer, err
}
