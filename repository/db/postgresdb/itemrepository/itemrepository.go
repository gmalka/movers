package itemrepository

import (
	"context"
	"fmt"

	"github.com/gmalka/movers/model"
	"github.com/jmoiron/sqlx"
)

type itemService struct {
	db *sqlx.DB
}

func NewItemService(db *sqlx.DB) *itemService {
	return &itemService{
		db: db,
	}
}

func (i *itemService) CreateItem(ctx context.Context, item model.Item) error {
	_, err := i.db.ExecContext(ctx, "INSERT INTO items(name,maxweight,minweight,maxprice,minprice) VALUES($1,$2,$3)", item.Name, item.MaxWeight, item.MinWeight)
	if err != nil {
		return fmt.Errorf("cant create item %s: %v", item.Name, err)
	}

	return nil
}

func (i *itemService) GetItemCount(ctx context.Context) (int, error) {
	var count int
	row := i.db.QueryRowContext(ctx, "SELECT COUNT(id) FROM items")
	if row.Err() != nil {
		return 0, fmt.Errorf("cant get count of items: %v", row.Err())
	}

	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("cant get count of items: %v", err)
	}

	return count, nil
}

func (i *itemService) GetItem(ctx context.Context, id int) (model.Item, error) {
	item := model.Item{}
	row := i.db.QueryRowContext(ctx, "SELECT name,maxweight,minweight FROM items WHERE id=$1", id)
	if row.Err() != nil {
		return model.Item{}, fmt.Errorf("cant get item %d: %v", id, row.Err())
	}

	err := row.Scan(&item.Name, &item.MaxWeight, &item.MinWeight)
	if err != nil {
		return model.Item{}, fmt.Errorf("cant get item %d: %v", id, err)
	}

	return item, nil
}
