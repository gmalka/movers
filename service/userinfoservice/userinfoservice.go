package userinfoservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/gmalka/movers/model"
)

type userInfoService struct {
	customers customerStorage
	workers   workerStorage
}

func NewuserInfoService(customers customerStorage, workers workerStorage) *userInfoService {
	return &userInfoService{
		customers: customers,
		workers:   workers,
	}
}

func (u *userInfoService) NewCustomer(ctx context.Context, customer model.CustomerInfo) error {
	if customer.Money < 10000 || customer.Money > 100000 {
		return errors.New("incorrect input data")
	}
	err := u.customers.CreateCustomer(ctx, customer)
	if err != nil {
		return fmt.Errorf("cant create new customer: %v", err)
	}

	return nil
}

func (u *userInfoService) GetCustomer(ctx context.Context, name string) (model.CustomerInfo, error) {
	return u.customers.GetCustomer(ctx, name)
}

func (u *userInfoService) DeleteCustomer(ctx context.Context, name string)  error {
	err := u.workers.UnchooseAll(ctx)
	if err != nil {
		return fmt.Errorf("cant delete customer: %v", err)
	}
	return u.customers.DeleteCustomer(ctx, name)
}

func (u *userInfoService) NewWorker(ctx context.Context, worker model.WorkerInfo) error {
	if worker.Salary < 10000 || worker.Salary > 30000 || worker.CarryWeight < 5 || worker.CarryWeight > 30 || worker.Fatigue < 0 || worker.Fatigue > 100 {
		return errors.New("incorrect input data")
	}
	err := u.workers.CreateWorker(ctx, worker)
	if err != nil {
		return fmt.Errorf("cant create new worker: %v", err)
	}

	return nil
}

func (u *userInfoService) GetWorker(ctx context.Context, name string) (model.WorkerInfo, error) {
	return u.workers.GetWorker(ctx, name)
}

func (u *userInfoService) GetWorkers(ctx context.Context) ([]model.WorkerInfo, error) {
	return u.workers.GetWorkers(ctx)
}

func (u *userInfoService) DeleteWorker(ctx context.Context, name string) error {
	return u.workers.DeleteWorker(ctx, name)
}

func (u *userInfoService) GetChoosenWorkers(ctx context.Context) ([]model.WorkerInfo, error) {
	return u.workers.GetChoosenWorkers(ctx)
}

func (u *userInfoService) RechooseWorkers(ctx context.Context, workers []string) error {
	return u.workers.RechooseWorkers(ctx, workers)
}

// <----------------INTERFACES---------------->

type customerStorage interface {
	CreateCustomer(ctx context.Context, customer model.CustomerInfo) error
	GetCustomer(ctx context.Context, name string) (model.CustomerInfo, error)
	DeleteCustomer(ctx context.Context, name string)  error 
}

type workerStorage interface {
	CreateWorker(ctx context.Context, worker model.WorkerInfo) error
	GetWorker(ctx context.Context, name string) (model.WorkerInfo, error)
	GetWorkers(ctx context.Context) ([]model.WorkerInfo, error)
	GetChoosenWorkers(ctx context.Context) ([]model.WorkerInfo, error)
	RechooseWorkers(ctx context.Context, workers []string) error
	UnchooseAll(ctx context.Context) error
	DeleteWorker(ctx context.Context, name string) error
}
