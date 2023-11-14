package userinfoservice

import (
	"fmt"

	"github.com/gmalka/movers/model"
)

type userInfoService struct {
	customers customerStorage
	workers workerStorage
}

func NewuserInfoService(customers customerStorage, workers workerStorage) userInfoService {
	return userInfoService{
		customers: customers,
		workers: workers,
	}
}

func (u *userInfoService) NewCustomer(customer model.CustomerInfo) error {
	err := u.customers.CreateCustomer(customer)
	if err != nil {
		return fmt.Errorf("cant create new customer: %v", err)
	}

	return nil
}

func (u *userInfoService) GetCustomer(name string) (model.CustomerInfo, error) {
	return u.customers.GetCustomer(name)
}

func (u *userInfoService) NewWorker(worker model.WorkerInfo) error {
	err := u.workers.CreateWorker(worker)
	if err != nil {
		return fmt.Errorf("cant create new worker: %v", err)
	}

	return nil
}

func (u *userInfoService) GetWorker(name string) (model.WorkerInfo, error) {
	return u.workers.GetWorker(name)
}

func (u *userInfoService) GetWorkers() ([]model.WorkerInfo, error) {
	return u.workers.GetWorkers()
}

// <----------------INTERFACES---------------->

type customerStorage interface {
	CreateCustomer(customer model.CustomerInfo) error
	GetCustomer(name string) (model.CustomerInfo, error)
}

type workerStorage interface {
	CreateWorker(worker model.WorkerInfo) error
	GetWorker(name string) (model.WorkerInfo, error)
	GetWorkers() ([]model.WorkerInfo, error)
}