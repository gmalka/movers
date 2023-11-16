package workservice

import (
	"context"
	"errors"
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

func (w *workService) MakeTaskForCustomer(ctx context.Context, customername string) error {
	task, err := w.tasker.GetFirstTask(ctx)
	if err != nil {
		return fmt.Errorf("cant calculate work: %v", err)
	}

	workers, err := w.workers.GetChoosenWorkers(ctx)
	if err != nil {
		return fmt.Errorf("cant calculate work: %v", err)
	}

	customer, err := w.customers.GetCustomer(ctx, customername)
	if err != nil {
		return fmt.Errorf("cant calculate work: %v", err)
	}

	if customer.Lost {
		return errors.New("user already loose")
	}

	salarySum := 0
	liftingCapacity := 0
	for _, v := range workers {
		salarySum += v.Salary
		if v.Drunk == 1 {
			liftingCapacity += int(float64(v.CarryWeight) * ((100 - float64(v.Fatigue))/100))
		} else {
			liftingCapacity += int(float64(v.CarryWeight) * ((100 - float64(v.Fatigue))/100) * (float64(50) / 100))
		}
	}

	if liftingCapacity < task.Weight {
		customer.Lost = true
		err = w.customers.UpdateCustomer(ctx, customer)
		if err != nil {
			return fmt.Errorf("cant calculate work: %v", err)
		}
		return fmt.Errorf("workers have not enough lifting capacity: have %v, want %v", liftingCapacity, task.Weight)
	}

	if customer.Money < salarySum {
		customer.Lost = true
		err = w.customers.UpdateCustomer(ctx, customer)
		if err != nil {
			return fmt.Errorf("cant calculate work: %v", err)
		}
		return fmt.Errorf("user %s has not enought money: have %v, want %v", customer.Name, liftingCapacity, task.Weight)
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
			return fmt.Errorf("cant calculate work3: %v", err)
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
	GetChoosenWorkers(ctx context.Context) ([]model.WorkerInfo, error)
	UpdateWorker(ctx context.Context, worker model.WorkerInfo) error
}

type taskFinisher interface {
	FinishTask(ctx context.Context, workers []string, task model.Task) error
	GetFirstTask(ctx context.Context) (model.Task, error)
}
