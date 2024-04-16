package services

import (
	"context"
	"errors"
	"myapp/internal/models"
	"myapp/internal/repositories"
)

type FinanceService struct {
	repo repositories.FinanceRepositories
}

func NewFinanceService(repo repositories.FinanceRepositories) *FinanceService {
	return &FinanceService{repo: repo}
}

func validateFinanceData(finance models.Finance) error {
	if finance.Currency == "" {
		return errors.New("currency must be provided")
	}
	if finance.Price == 0 {
		return errors.New("the value cannot be zero")
	}
	if finance.Folder_id == 0 {
		return errors.New("the value folder_id cannot be zero")
	}

	return nil
}

func (s *FinanceService) CreateFinance(ctx context.Context, finance models.Finance) error {
	if err := validateFinanceData(finance); err != nil {
		return nil
	}
	return s.repo.CreateFinance(ctx, finance)
}

func (s *FinanceService) DeleteFinance(id_to_del int) error {
	return s.repo.DeleteFinance(context.Background(), id_to_del)
}

func (s *FinanceService) FetchFinance(start, end, month int) ([][]string, error) {
	return s.repo.FetchFinance(context.Background(), start, end, month)
}
