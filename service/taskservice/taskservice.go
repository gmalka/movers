package taskservice

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gmalka/movers/model"
)

type taskService struct {
	itemcash map[int]model.Item
	it       itemStore
	ts       taskStore
}

func NewTaskService(it itemStore, ts taskStore) *taskService {
	return &taskService{
		it: it,
		ts: ts,
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

	t.itemcash = make(map[int]model.Item, 5)
	rand.Seed(time.Now().Unix())
	tasks := make([]model.Task, tocreate)
	for i := 0; i < tocreate; i++ {
		task := model.Task{}
		id := rand.Intn(count)

		var item model.Item
		if val, ok := t.itemcash[id]; ok {
			item = val
		} else {
			item, err = t.it.GetItem(id)
			if err != nil {
				return fmt.Errorf("cant generate task; %v", err)
			}
			t.itemcash[id] = item
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

func (t *taskService) FinishTask(taskId int) error {
	return t.ts.DeleteTask(taskId)
}

// <----------------INTERFACES---------------->

type itemStore interface {
	CreateItem(item model.Item) error
	GetItemCount() (int, error)
	GetItem(id int) (model.Item, error)
}

type taskStore interface {
	CreateTasks(tasks []model.Task) error
	GetTasks() ([]model.Task, error)
	DeleteTask(taskId int) error
	DeleteTasks() error
}
