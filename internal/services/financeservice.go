package services

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/models"
	"myapp/internal/repositories"
	"time"

	"github.com/patrickmn/go-cache"
)

type FinanceService struct {
	repo  repositories.FinanceRepositories
	cache *cache.Cache
}

func NewFinanceService(repo repositories.FinanceRepositories) *FinanceService {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &FinanceService{
		repo:  repo,
		cache: c,
	}
}

func (s *FinanceService) AddCacheKeyForUser(user_id int, key string) {
	userCacheKey := fmt.Sprintf("finance-user-%d-keys", user_id)
	var keys []string
	if cachedKeys, found := s.cache.Get(userCacheKey); found {
		keys = cachedKeys.([]string)
	}
	keys = append(keys, key)
	s.cache.Set(userCacheKey, keys, cache.DefaultExpiration)
}

func (s *FinanceService) InvalidataUserCache(user_id int) {
	userCacheKey := fmt.Sprintf("finance-user-%d-keys", user_id)
	if keys, found := s.cache.Get(userCacheKey); found {
		for _, key := range keys.([]string) {
			s.cache.Delete(key)
		}
		s.cache.Set(userCacheKey, []string{}, cache.DefaultExpiration)
	}
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
		return err
	}

	finId, err := s.repo.CreateFinance(ctx, finance)
	if err != nil {
		return err
	}

	userId, err := s.repo.GetUserByFinance(finId)
	if err != nil {
		return err
	}

	// Инвалидация кэша для одного пользователя
	s.InvalidataUserCache(userId)
	return nil

}

func (s *FinanceService) DeleteFinance(ctx context.Context, id_to_del int) error {
	userId, err := s.repo.GetUserByFinance(id_to_del)
	if err != nil {
		return err
	}

	err = s.repo.DeleteFinance(ctx, id_to_del)
	if err != nil {
		return err
	}

	// Инвалидация кэша для одного пользователя
	s.InvalidataUserCache(userId)
	return nil
}

func (s *FinanceService) FetchFinance(user_id, start, end int) ([]models.Finance, error) {
	key := fmt.Sprintf("finance-%d-%d-%d", user_id, start, end)

	if cachedData, found := s.cache.Get(key); found {
		return cachedData.([]models.Finance), nil
	}

	data, err := s.repo.FetchFinance(context.Background(), user_id, start, end)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(user_id, key)
	return data, nil
}

func (s *FinanceService) FetchFinanceIncome(user_id, start, end int, yearMonth string) ([]models.Finance, error) {
	key := fmt.Sprintf("finance-income-%d-%d-%d-%q", user_id, start, end, yearMonth)

	if cachedData, found := s.cache.Get(key); found {
		return cachedData.([]models.Finance), nil
	}

	data, err := s.repo.FetchFinanceIncome(context.Background(), user_id, start, end, yearMonth)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(user_id, key)
	return data, nil
}

func (s *FinanceService) FetchFinanceExpense(user_id, start, end int, yearMonth string) ([]models.Finance, error) {
	key := fmt.Sprintf("finance-expense-%d-%d-%d-%q", user_id, start, end, yearMonth)

	if cachedData, found := s.cache.Get(key); found {
		return cachedData.([]models.Finance), nil
	}

	data, err := s.repo.FetchFinanceExpense(context.Background(), user_id, start, end, yearMonth)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(user_id, key)
	return data, nil
}
