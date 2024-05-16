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

type FolderService struct {
	repo  repositories.FolderRepositories
	cache *cache.Cache
}

func NewFolderService(repo repositories.FolderRepositories) *FolderService {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &FolderService{
		repo:  repo,
		cache: c,
	}
}

func (s *FolderService) AddCacheKeyForUser(user_id int, key string) {
	userCacheKey := fmt.Sprintf("folder-user-%d-keys", user_id)
	var keys []string
	if cacheKey, found := s.cache.Get(userCacheKey); found {
		keys = cacheKey.([]string)
	}
	keys = append(keys, key)
	s.cache.Set(userCacheKey, keys, cache.DefaultExpiration)
}

func (s *FolderService) InvalidataUserCache(user_id int) {
	userCacheKey := fmt.Sprintf("folder-user-%d-keys", user_id)
	if keys, found := s.cache.Get(userCacheKey); found {
		for _, key := range keys.([]string) {
			s.cache.Delete(key)
		}
		s.cache.Set(userCacheKey, []string{}, cache.DefaultExpiration)
	}
}

func validateFolderData(folder models.Folder) error {
	if folder.Name == "" || folder.Type == "" || folder.Image == "" {
		return errors.New("name, type and image must be provided")
	}
	if len(folder.Name) > 25 {
		return errors.New("the length of the name cannot exceed 25 characters")
	}
	if len(folder.Name) < 3 {
		return errors.New("the length of the name cannot be shorter than 3 characters")
	}
	return nil
}

func (s *FolderService) CreateFolder(ctx context.Context, folder models.Folder) error {
	if err := validateFolderData(folder); err != nil {
		return err
	}

	err := s.repo.CreateFolder(ctx, folder)
	if err != nil {
		return err
	}

	s.InvalidataUserCache(*folder.User_id)
	return nil
}

func (s *FolderService) DeleteFolder(id_to_del int) error {
	userId, err := s.repo.GetUserByFolder(id_to_del)
	if err != nil {
		return err
	}

	err = s.repo.DeleteFolder(context.Background(), id_to_del)
	if err != nil {
		return err
	}

	s.InvalidataUserCache(userId)
	return nil
}

func (s *FolderService) FetchFolder(start, end int, id_user int, type_folder string) ([]models.Folder, error) {
	key := fmt.Sprintf("folder-%d-%d-%d-%q", id_user, start, end, type_folder)

	if cachedData, found := s.cache.Get(key); found {
		return cachedData.([]models.Folder), nil
	}

	data, err := s.repo.FetchFolder(context.Background(), start, end, id_user, type_folder)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(id_user, key)
	return data, nil
}

func (s *FolderService) UpdateFolder(ctx context.Context, folder models.Folder) error {
	if err := validateFolderData(folder); err != nil {
		return err
	}

	err := s.repo.UpdateFolder(ctx, folder)
	if err != nil {
		return err
	}

	s.InvalidataUserCache(*folder.User_id)
	return nil
}

func (s *FolderService) FetchFolderById(id_folder, user_id int) ([]models.Folder, error) {
	key := fmt.Sprintf("folder-%d-%d", user_id, id_folder)

	if cachedData, found := s.cache.Get(key); found {
		return cachedData.([]models.Folder), nil
	}

	data, err := s.repo.FetchFolderById(context.Background(), id_folder, user_id)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(user_id, key)
	return data, nil
}
