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

type TaskServices struct {
	repo  repositories.TaskRepositories
	cache *cache.Cache
}

func NewTaskServices(repo repositories.TaskRepositories) *TaskServices {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &TaskServices{
		repo:  repo,
		cache: c,
	}
}

func (s *TaskServices) AddCacheKeyForUser(user_id int, key string) {
	userCacheKey := fmt.Sprintf("task-user-%d-keys", user_id)
	var keys []string
	if cacheKey, found := s.cache.Get(userCacheKey); found {
		keys = cacheKey.([]string)
	}
	keys = append(keys, key)
	s.cache.Set(userCacheKey, keys, cache.DefaultExpiration)
}

func (s *TaskServices) InvalidataUserCache(user_id int) {
	userCacheKey := fmt.Sprintf("task-user-%d-keys", user_id)
	if keys, found := s.cache.Get(userCacheKey); found {
		for _, key := range keys.([]string) {
			s.cache.Delete(key)
		}
		s.cache.Set(userCacheKey, []string{}, cache.DefaultExpiration)
	}
}

func validateTaskData(task models.Task) error {
	if task.Text == "" {
		return errors.New("text must be provided")
	}
	if task.Folder_id == 0 {
		return errors.New("the value folder_id cannot be zero")
	}

	return nil
}

func (s *TaskServices) CreateTask(ctx context.Context, task models.Task) error {
	if err := validateTaskData(task); err != nil {
		return err
	}

	taskId, err := s.repo.CreateTask(ctx, task)
	if err != nil {
		return err
	}

	userID, err := s.repo.GetUserByTask(taskId)
	if err != nil {
		return err
	}
	s.InvalidataUserCache(userID)
	return nil
}

func (s *TaskServices) DeleteTask(id_to_del int) error {
	userID, err := s.repo.GetUserByTask(id_to_del)
	if err != nil {
		return err
	}

	err = s.repo.DeleteTask(context.Background(), id_to_del)
	if err != nil {
		return err
	}

	s.InvalidataUserCache(userID)
	return nil
}

func (s *TaskServices) FetchTask(user_id, start, end int, folder_id ...int) ([]models.Task, error) {
	var key string
	if len(folder_id) > 0 {
		key = fmt.Sprintf("task-%d-%d-%d-%d", user_id, start, end, folder_id[0])
	} else {
		key = fmt.Sprintf("task-%d-%d-%d", user_id, start, end)
	}

	if cachedData, found := s.cache.Get(key); found {
		return cachedData.([]models.Task), nil
	}

	data, err := s.repo.FetchTask(context.Background(), user_id, start, end, folder_id...)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(user_id, key)
	return data, nil
}

func (s *TaskServices) UpdateTask(ctx context.Context, task models.Task) error {
	userID, err := s.repo.GetUserByTask(task.Id)
	if err != nil {
		return err
	}

	err = s.repo.UpdateTask(ctx, task)
	if err != nil {
		return err
	}

	s.InvalidataUserCache(userID)
	return nil
}

func (s *TaskServices) CountTask(user_id int) (int, error) {
	key := fmt.Sprintf("task-count-%d", user_id)
	if cachedData, found := s.cache.Get(key); found {
		return cachedData.(int), nil
	}

	data, err := s.repo.CountTask(context.Background(), user_id)
	if err != nil {
		return 0, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(user_id, key)
	return data, nil
}

func (s *TaskServices) CountTaskFavourites(user_id int) (int, error) {
	key := fmt.Sprintf("task-count-favourites-%d", user_id)
	if cachedData, found := s.cache.Get(key); found {
		return cachedData.(int), nil
	}

	data, err := s.repo.CountTaskFavourites(context.Background(), user_id)
	if err != nil {
		return 0, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(user_id, key)
	return data, nil
}

func (s *TaskServices) FetchTaskFavourites(user_id, start, end int, folder_id ...int) ([]models.Task, error) {
	var key string
	if len(folder_id) > 0 {
		key = fmt.Sprintf("task-favourites-%d-%d-%d-%d", user_id, start, end, folder_id[0])
	} else {
		key = fmt.Sprintf("task-favourites-%d-%d-%d", user_id, start, end)
	}

	if cachedData, found := s.cache.Get(key); found {
		return cachedData.([]models.Task), nil
	}

	data, err := s.repo.FetchTaskFavourites(context.Background(), user_id, start, end, folder_id...)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(user_id, key)
	return data, nil
}
