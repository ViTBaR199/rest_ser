package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/models"
	"time"

	_ "github.com/lib/pq"
)

type FinanceRepositories interface {
	CreateFinance(ctx context.Context, finance models.Finance) (int, error)
	DeleteFinance(ctx context.Context, id_to_del int) error
	FetchFinance(ctx context.Context, start, end, month int) ([]models.Finance, error)
	FetchFinanceIncome(ctx context.Context, user_id, start, end int, yearMonth string) ([]models.Finance, error)
	FetchFinanceExpense(ctx context.Context, user_id, start, end int, yearMonth string) ([]models.Finance, error)
	GetUserByFinance(financeID int) (int, error)
}

type financeRepositories struct {
	db *sql.DB
}

func NewFinanceRepositories(db *sql.DB) FinanceRepositories {
	return &financeRepositories{db: db}
}

func (r *financeRepositories) CreateFinance(ctx context.Context, finance models.Finance) (int, error) {
	var newId int
	err := r.db.QueryRowContext(ctx, "SELECT create_new_finance($1, $2, $3, $4)", finance.Price, finance.Currency, finance.Folder_id, finance.Date).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (r *financeRepositories) DeleteFinance(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_finance($1)", id_to_del)
	return err
}

func (r *financeRepositories) FetchFinance(ctx context.Context, user_id, start, end int) ([]models.Finance, error) {
	var result []models.Finance

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM fetch_finance($1, $2, $3)", user_id, start, end)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_finance: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Finance
		var date_d time.Time
		if err := rows.Scan(&f.Id, &f.CategoryName, &f.CategoryPhoto, &f.CategoryColor, &f.Price, &date_d, &f.Currency, &f.Folder_id); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		f.Date = time.Date(date_d.Year(), date_d.Month(), date_d.Day(), 0, 0, 0, 0, date_d.Location()).Format("2006.01.02")
		result = append(result, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return result, nil
}

func (r *financeRepositories) FetchFinanceIncome(ctx context.Context, user_id, start, end int, yearMonth string) ([]models.Finance, error) {
	var result []models.Finance

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM fetch_finance_income($1, $2, $3, $4)", user_id, start, end, yearMonth)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_finance_income: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Finance
		var date_d time.Time
		if err := rows.Scan(&f.Id, &f.CategoryName, &f.CategoryPhoto, &f.CategoryColor, &f.Price, &date_d, &f.Currency, &f.Folder_id); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		f.Date = time.Date(date_d.Year(), date_d.Month(), date_d.Day(), 0, 0, 0, 0, date_d.Location()).Format("2006.01.02")
		result = append(result, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return result, nil
}

func (r *financeRepositories) FetchFinanceExpense(ctx context.Context, user_id, start, end int, yearMonth string) ([]models.Finance, error) {
	var result []models.Finance

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM fetch_finance_expense($1, $2, $3, $4)", user_id, start, end, yearMonth)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_finance_expense: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Finance
		var date_d time.Time
		if err := rows.Scan(&f.Id, &f.CategoryName, &f.CategoryPhoto, &f.CategoryColor, &f.Price, &date_d, &f.Currency, &f.Folder_id); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		f.Date = time.Date(date_d.Year(), date_d.Month(), date_d.Day(), 0, 0, 0, 0, date_d.Location()).Format("2006.01.02")
		result = append(result, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return result, nil
}

func (r *financeRepositories) GetUserByFinance(financeID int) (int, error) {
	var userId int
	err := r.db.QueryRow("SELECT get_user_by_finance($1)", financeID).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
