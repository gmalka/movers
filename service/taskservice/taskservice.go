package taskservice

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gmalka/movers/model"
)

type taskService struct {
	it itemStore
	ts taskStore
	dt doneTasksStore
}

func NewTaskService(it itemStore, ts taskStore, dt doneTasksStore) taskService {
	return taskService{
		it: it,
		ts: ts,
		dt: dt,
	}
}

func (t *taskService) CreateItem(ctx context.Context, item model.Item) error {
	return t.it.CreateItem(ctx, item)
}

func (t *taskService) GenerateTasks(ctx context.Context, tocreate int) error {
	count, err := t.it.GetItemCount(ctx)
	if err != nil {
		return fmt.Errorf("cant generate tasks; %v", err)
	}

	itemcash := make(map[int]model.Item, 5)
	rand.Seed(time.Now().Unix())
	tasks := make([]model.Task, tocreate)
	for i := 0; i < tocreate; i++ {
		task := model.Task{}
		id := rand.Intn(count)

		var item model.Item
		if val, ok := itemcash[id]; ok {
			item = val
		} else {
			item, err = t.it.GetItem(ctx, id)
			if err != nil {
				return fmt.Errorf("cant generate task; %v", err)
			}
			itemcash[id] = item
		}

		task.ItemName = item.Name
		task.Weight = rand.Intn(item.MaxWeight) + item.MinWeight

		tasks[i] = task
	}

	err = t.ts.CreateTasks(ctx, tasks)
	if err != nil {
		return fmt.Errorf("cant create tasks: %v", err)
	}

	return nil
}

func (t *taskService) GetFirstTask(ctx context.Context) (model.Task, error) {
	return t.ts.GetFirstTask(ctx)
}

func (t *taskService) GetTasks(ctx context.Context) ([]model.Task, error) {
	return t.ts.GetTasks(ctx)
}

func (t *taskService) GetWorkerTasks(ctx context.Context, name string) ([]model.Task, error) {
	return t.dt.GetWorkerTasks(ctx, name)
}

func (t *taskService) FinishTask(ctx context.Context, workers []string, task model.Task) error {
	err := t.ts.DeleteTask(ctx, task.TaskId)
	if err != nil {
		return fmt.Errorf("cant finish task: %v", err)
	}

	err = t.dt.CompleteTask(ctx, workers, task)
	if err != nil {
		return fmt.Errorf("cant finish task: %v", err)
	}

	return nil
}

// <----------------INTERFACES---------------->

type itemStore interface {
	CreateItem(ctx context.Context, item model.Item) error
	GetItemCount(ctx context.Context) (int, error)
	GetItem(ctx context.Context, id int) (model.Item, error)
}

type doneTasksStore interface {
	CompleteTask(ctx context.Context, workers []string, task model.Task) error
	GetWorkerTasks(ctx context.Context, name string) ([]model.Task, error)
}

type taskStore interface {
	CreateTasks(ctx context.Context, tasks []model.Task) error
	GetFirstTask(ctx context.Context) (model.Task, error)
	GetTasks(ctx context.Context) ([]model.Task, error)
	DeleteTask(ctx context.Context, taskId int) error
}
