package workservice

import (
	"context"
	"fmt"

	"github.com/gmalka/movers/model"
)

type workService struct {
	customers customerGetter
	workers   workerGetter
	tasker    taskFinisher
}

func NewWorkService(customers customerGetter, workers workerGetter, tasker taskFinisher) *workService {
	return &workService{
		customers: customers,
		workers:   workers,
		tasker:    tasker,
	}
}

func (w *workService) CalculateWork(ctx context.Context, customername string, workers []model.WorkerInfo, task model.Task) error {
	customer, err := w.customers.GetCustomer(ctx, customername)
	if err != nil {
		return fmt.Errorf("cant calculate work: %v", err)
	}

	salarySum := 0
	liftingCapacity := 0
	for _, v := range workers {
		salarySum += v.Salary
		liftingCapacity += v.CarryWeight * (100 - v.Fatigue/100) * (v.Drunk * 100)
	}

	if customer.Money < salarySum {
		return fmt.Errorf("user %s has not enought money", customer.Name)
	}

	customer.Money -= salarySum
	err = w.customers.UpdateCustomer(ctx, customer)
	if err != nil {
		return fmt.Errorf("cant calculate work: %v", err)
	}

	tasknames := make([]string, len(workers))
	for i := 0; i < len(workers); i++ {
		if workers[i].Fatigue < 80 {
			workers[i].Fatigue += 20
		} else {
			workers[i].Fatigue = 100
		}

		err = w.workers.UpdateWorker(ctx, workers[i])
		if err != nil {
			return fmt.Errorf("cant calculate work: %v", err)
		}

		tasknames[i] = workers[i].Name
	}

	err = w.tasker.FinishTask(ctx, tasknames, task)
	if err != nil {
		return fmt.Errorf("cant calculate work: %v", err)
	}

	return nil
}

// <----------------INTERFACES---------------->

type customerGetter interface {
	GetCustomer(ctx context.Context, name string) (model.CustomerInfo, error)
	UpdateCustomer(ctx context.Context, customer model.CustomerInfo) error
}

type workerGetter interface {
	GetWorker(ctx context.Context, name string) (model.WorkerInfo, error)
	UpdateWorker(ctx context.Context, worker model.WorkerInfo) error
}

type taskFinisher interface {
	FinishTask(ctx context.Context, workers []string, task model.Task) error
}
