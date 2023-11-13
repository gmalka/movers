package moverrepository

import (
	"fmt"

	"github.com/gmalka/movers/model"
	"github.com/jmoiron/sqlx"
)

type MoverRepository struct {
	db *sqlx.DB
}

func NewMoverRepository(db *sqlx.DB) *MoverRepository {
	return &MoverRepository{
		db: db,
	}
}

func (m *MoverRepository) CreateMover(mover model.MoverInfo) error {
	_, err := m.db.Exec("INSERT INTO movers VALUES($1,$2,$3,$4,$5)", mover.Name, mover.Fatigue, mover.Salary, mover.CarryWeight, mover.Drunk)
	if err != nil {
		return fmt.Errorf("cant insert mover %s: %v", mover.Name, err)
	}

	return nil
}

func (m *MoverRepository) GetMover(name string) (model.MoverInfo, error) {
	row := m.db.QueryRow("SELECT * FROM movers WHERE name=$1", name)
	if row.Err() != nil {
		return model.MoverInfo{}, fmt.Errorf("cant find mover %s: %v", name, row.Err())
	}

	mover := model.MoverInfo{}
	err := row.Scan(&mover.Name, &mover.Fatigue, &mover.Salary, &mover.CarryWeight, &mover.Drunk)
	if err != nil {
		return model.MoverInfo{}, fmt.Errorf("cant scan mover %s: %v", name, row.Err())
	}

	return mover, nil
}

func (m *MoverRepository) UpdateMover(mover model.MoverInfo) error {
	_, err := m.db.Exec("UPDATE movers SET name = $1, fatigue = $2, salary = $3, carryweight = $4, drunk = $5", mover.Name, mover.Fatigue, mover.Salary, mover.CarryWeight, mover.Drunk)
	if err != nil {
		return fmt.Errorf("cant update mover %s: %v", mover.Name, err)
	}

	return nil
}