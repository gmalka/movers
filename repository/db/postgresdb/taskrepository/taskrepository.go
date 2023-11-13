package taskrepository

import (
	"fmt"

	"github.com/gmalka/movers/model"
	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return TaskRepository{
		db: db,
	}
}

func (t *TaskRepository) CreateTasks(tasks []model.Task) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("cant create tasks transaction: %v", err)
	}
	defer tx.Rollback()

	for _, v := range tasks {
		_, err = tx.Exec("INSERT INTO tasks(itemname,weight) VALUES($1,$2)", v.ItemName, v.Weight)
		if err != nil {
			return fmt.Errorf("cant insert task: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit tasks transaction: %v", err)
	}
	return nil
}

func (t *TaskRepository) GetTasks() ([]model.Task, error) {
	tasks := make([]model.Task, 0, 10)
	rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("cant select tasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := model.Task{}
		err = rows.Scan(&task.TaskId, &task.ItemName, &task.Weight)
		if err != nil {
			return nil, fmt.Errorf("cant scan task: %v", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *TaskRepository) DeleteTask(taskId int) error {
	_, err := t.db.Exec("DELETE FROM tasks WHERE id=$1", taskId)
	if err != nil {
		return fmt.Errorf("cant delete task %d: %v", taskId, err)
	}

	return nil
}

func (t *TaskRepository) DeleteTasks() error {
	_, err := t.db.Exec("DELETE FROM tasks")
	if err != nil {
		return fmt.Errorf("cant delete tasks: %v", err)
	}

	return nil
}