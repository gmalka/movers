package taskservice

import (
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

func (t *taskService) CreateItem(item model.Item) error {
	return t.it.CreateItem(item)
}

func (t *taskService) GenerateTasks(tocreate int) error {
	count, err := t.it.GetItemCount()
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
			item, err = t.it.GetItem(id)
			if err != nil {
				return fmt.Errorf("cant generate task; %v", err)
			}
			itemcash[id] = item
		}

		task.ItemName = item.Name
		task.Weight = rand.Intn(item.MaxWeight) + item.MinWeight

		tasks[i] = task
	}

	err = t.ts.CreateTasks(tasks)
	if err != nil {
		return fmt.Errorf("cant create tasks: %v", err)
	}

	return nil
}

func (t *taskService) GetTasks() ([]model.Task, error) {
	return t.ts.GetTasks()
}

func (t *taskService) GetWorkerTasks(name string) ([]model.Task, error) {
	return t.dt.GetWorkerTasks(name)
}

func (t *taskService) FinishTask(workers []string, task model.Task) error {
	err := t.ts.DeleteTask(task.TaskId)
	if err != nil {
		return fmt.Errorf("cant finish task: %v", err)
	}

	err = t.dt.CompleteTask(workers, task)
	if err != nil {
		return fmt.Errorf("cant finish task: %v", err)
	}

	return nil
}

// <----------------INTERFACES---------------->

type itemStore interface {
	CreateItem(item model.Item) error
	GetItemCount() (int, error)
	GetItem(id int) (model.Item, error)
}

type doneTasksStore interface {
	CompleteTask(workers []string, task model.Task) error
	GetWorkerTasks(name string) ([]model.Task, error)
}

type taskStore interface {
	CreateTasks(tasks []model.Task) error
	GetTasks() ([]model.Task, error)
	DeleteTask(taskId int) error
}
