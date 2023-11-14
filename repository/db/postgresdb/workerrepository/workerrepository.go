package workerrepository

import (
	"context"
	"fmt"

	"github.com/gmalka/movers/model"
	"github.com/jmoiron/sqlx"
)

type WorkerRepository struct {
	db *sqlx.DB
}

func NewWorkerRepository(db *sqlx.DB) *WorkerRepository {
	return &WorkerRepository{
		db: db,
	}
}

func (w *WorkerRepository) CreateWorker(ctx context.Context, worker model.WorkerInfo) error {
	_, err := w.db.ExecContext(ctx, "INSERT INTO workers VALUES($1,$2,$3,$4,$5)", worker.Name, worker.Fatigue, worker.Salary, worker.CarryWeight, worker.Drunk)
	if err != nil {
		return fmt.Errorf("cant insert worker %s: %v", worker.Name, err)
	}

	return nil
}

func (w *WorkerRepository) GetWorker(ctx context.Context, name string) (model.WorkerInfo, error) {
	row := w.db.QueryRowContext(ctx, "SELECT * FROM workers WHERE name=$1", name)
	if row.Err() != nil {
		return model.WorkerInfo{}, fmt.Errorf("cant find worker %s: %v", name, row.Err())
	}

	worker := model.WorkerInfo{}
	err := row.Scan(&worker.Name, &worker.Fatigue, &worker.Salary, &worker.CarryWeight, &worker.Drunk)
	if err != nil {
		return model.WorkerInfo{}, fmt.Errorf("cant scan worker %s: %v", name, row.Err())
	}

	return worker, nil
}

func (w *WorkerRepository) GetWorkers(ctx context.Context, ) ([]model.WorkerInfo, error) {
	workers := make([]model.WorkerInfo, 0, 10)
	rows, err := w.db.QueryContext(ctx, "SELECT * FROM workers")
	if err != nil {
		return nil, fmt.Errorf("cant select workers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		worker := model.WorkerInfo{}
		err := rows.Scan(&worker.Name, &worker.Fatigue, &worker.Salary, &worker.CarryWeight, &worker.Drunk)
		if err != nil {
			return nil, fmt.Errorf("cant scan task: %v", err)
		}

		workers = append(workers, worker)
	}

	return workers, nil
}

func (w *WorkerRepository) UpdateWorker(ctx context.Context, worker model.WorkerInfo) error {
	_, err := w.db.ExecContext(ctx, "UPDATE workers SET name = $1, fatigue = $2, salary = $3, carryweight = $4, drunk = $5", worker.Name, worker.Fatigue, worker.Salary, worker.CarryWeight, worker.Drunk)
	if err != nil {
		return fmt.Errorf("cant update worker %s: %v", worker.Name, err)
	}

	return nil
}
