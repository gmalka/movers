package workerrepository

import (
	"context"
	"fmt"

	"github.com/gmalka/movers/model"
	"github.com/jmoiron/sqlx"
)

type workerRepository struct {
	db *sqlx.DB
}

func NewWorkerRepository(db *sqlx.DB) *workerRepository {
	return &workerRepository{
		db: db,
	}
}

func (w *workerRepository) CreateWorker(ctx context.Context, worker model.WorkerInfo) error {
	_, err := w.db.ExecContext(ctx, "INSERT INTO workers VALUES($1,$2,$3,$4,$5,$6)", worker.Name, worker.Fatigue, worker.Salary, worker.CarryWeight, worker.Drunk, false)
	if err != nil {
		return fmt.Errorf("cant insert worker %s: %v", worker.Name, err)
	}

	return nil
}

func (w *workerRepository) GetWorker(ctx context.Context, name string) (model.WorkerInfo, error) {
	row := w.db.QueryRowContext(ctx, "SELECT * FROM workers WHERE name=$1", name)
	if row.Err() != nil {
		return model.WorkerInfo{}, fmt.Errorf("cant find worker %s: %v", name, row.Err())
	}

	worker := model.WorkerInfo{}
	err := row.Scan(&worker.Name, &worker.Fatigue, &worker.Salary, &worker.CarryWeight, &worker.Drunk, &worker.Choosen)
	if err != nil {
		return model.WorkerInfo{}, fmt.Errorf("cant scan worker %s: %v", name, row.Err())
	}

	return worker, nil
}

func (w *workerRepository) GetWorkers(ctx context.Context) ([]model.WorkerInfo, error) {
	workers := make([]model.WorkerInfo, 0, 10)
	rows, err := w.db.QueryContext(ctx, "SELECT * FROM workers")
	if err != nil {
		return nil, fmt.Errorf("cant select workers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		worker := model.WorkerInfo{}
		err := rows.Scan(&worker.Name, &worker.Fatigue, &worker.Salary, &worker.CarryWeight, &worker.Drunk, &worker.Choosen)
		if err != nil {
			return nil, fmt.Errorf("cant scan task: %v", err)
		}

		workers = append(workers, worker)
	}

	return workers, nil
}

func (w *workerRepository) GetChoosenWorkers(ctx context.Context) ([]model.WorkerInfo, error) {
	workers := make([]model.WorkerInfo, 0, 10)
	rows, err := w.db.QueryContext(ctx, "SELECT * FROM workers WHERE choosen = true")
	if err != nil {
		return nil, fmt.Errorf("cant select workers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		worker := model.WorkerInfo{}
		err := rows.Scan(&worker.Name, &worker.Fatigue, &worker.Salary, &worker.CarryWeight, &worker.Drunk, &worker.Choosen)
		if err != nil {
			return nil, fmt.Errorf("cant scan task: %v", err)
		}

		workers = append(workers, worker)
	}

	return workers, nil
}

func (w *workerRepository) ChooseWorkers(ctx context.Context, workers []model.WorkerInfo) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cant start transaction: %v", err)
	}
	defer tx.Rollback()

	for _, v := range workers {
		_, err := tx.ExecContext(ctx, "UPDATE workers SET choosen = $1 WHERE name = $2", true, v.Name)
		if err != nil {
			return fmt.Errorf("cant choose %s: %v", v.Name, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit tasks transaction: %v", err)
	}
	return nil
}

func (w *workerRepository) UnchooseWorkers(ctx context.Context, workers []model.WorkerInfo) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cant start transaction: %v", err)
	}
	defer tx.Rollback()

	for _, v := range workers {
		_, err := tx.ExecContext(ctx, "UPDATE workers SET choosen = $1 WHERE name = $2", false, v.Name)
		if err != nil {
			return fmt.Errorf("cant choose %s: %v", v.Name, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit tasks transaction: %v", err)
	}
	return nil
}

func (w *workerRepository) UnchooseAll(ctx context.Context) error {
	_, err := w.db.ExecContext(ctx, "UPDATE workers SET choosen = $1", false)
	if err != nil {
		return fmt.Errorf("cant unchoose all: %v", err)
	}

	return nil
}

func (w *workerRepository) UpdateWorker(ctx context.Context, worker model.WorkerInfo) error {
	_, err := w.db.ExecContext(ctx, "UPDATE workers SET name = $1, fatigue = $2, salary = $3, carryweight = $4, drunk = $5", worker.Name, worker.Fatigue, worker.Salary, worker.CarryWeight, worker.Drunk)
	if err != nil {
		return fmt.Errorf("cant update worker %s: %v", worker.Name, err)
	}

	return nil
}

func (w *workerRepository) DeleteWorker(ctx context.Context, name string) error {
	_, err := w.db.ExecContext(ctx, "DELETE FROM workers WHERE name=$1", name)
	if err != nil {
		return fmt.Errorf("cant delete worker %s: %v", name, err)
	}

	return nil
}
