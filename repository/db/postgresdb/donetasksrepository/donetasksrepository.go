package donetasksrepository

import (
	"context"
	"fmt"

	"github.com/gmalka/movers/model"
	"github.com/jmoiron/sqlx"
)

const pageLimit = 20

type doneTasksRepository struct {
	db *sqlx.DB
}

func NewDoneTasksRepository(db *sqlx.DB) *doneTasksRepository {
	return &doneTasksRepository{
		db: db,
	}
}

func (d *doneTasksRepository) CompleteTask(ctx context.Context, workers []string, task model.Task) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cant start transaction: %v", err)
	}
	defer tx.Rollback()

	for _, v := range workers {
		_, err := tx.ExecContext(ctx, "INSERT INTO completetasks(workername,itemname,weight) VALUES($1,$2,$3)", v, task.ItemName, task.Weight)
		if err != nil {
			return fmt.Errorf("cant add task %s to completetask task: %v", v, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit tasks transaction: %v", err)
	}
	return nil
}

func (d *doneTasksRepository) GetWorkerTasks(ctx context.Context, name string, page int) ([]model.Task, error) {
	page -= 1
	if page < 0 {
		page = 0
	}

	rows, err := d.db.QueryContext(ctx, "SELECT id,itemname,weight FROM completetasks WHERE workername=$1 LIMIT $2 OFFSET $3", name, pageLimit, page * pageLimit)
	if err != nil {
		return nil, fmt.Errorf("cant get worker tasks: %v", err)
	}

	tasks := make([]model.Task, 0, 10)
	for rows.Next() {
		task := model.Task{}
		err = rows.Scan(&task.TaskId, &task.ItemName, &task.Weight)
		if err != nil {
			return nil, fmt.Errorf("cant get worker tasks: %v", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
