package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/models"

	_ "github.com/lib/pq"
)

type FinanceRepositories interface {
	CreateFinance(ctx context.Context, finance models.Finance) error
	DeleteFinance(ctx context.Context, id_to_del int) error
	FetchFinance(ctx context.Context, start, end, month int) ([][]string, error)
}

type financeRepositories struct {
	db *sql.DB
}

func NewFinanceRepositories(db *sql.DB) FinanceRepositories {
	return &financeRepositories{db: db}
}

func (r *financeRepositories) CreateFinance(ctx context.Context, finance models.Finance) error {
	_, err := r.db.ExecContext(ctx, "SELECT create_new_finance($1, $2, $3, $4)", finance.Price, finance.Currency, finance.Folder_id, finance.Date)
	return err
}

func (r *financeRepositories) DeleteFinance(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_finance($1)", id_to_del)
	return err
}

func (r *financeRepositories) FetchFinance(ctx context.Context, start, end, month int) ([][]string, error) {
	var result [][]string

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM fetch_finance($1, $2, $3)", start, end, month)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_finance: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, price, folder_id int
		var currancy string
		var date string
		if err := rows.Scan(&id, &price, &date, &currancy, &folder_id); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		rowData := []string{fmt.Sprintf("%d", id), fmt.Sprintf("%d", price), date, currancy, fmt.Sprintf("%d", folder_id)}
		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return result, nil
}
