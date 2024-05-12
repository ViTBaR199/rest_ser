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

type NoteService struct {
	repo  repositories.NoteRepositories
	cache *cache.Cache
}

func NewNoteService(repo repositories.NoteRepositories) *NoteService {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &NoteService{
		repo:  repo,
		cache: c,
	}
}

func (s *NoteService) AddCacheKeyForUser(user_id int, key string) {
	userCacheKey := fmt.Sprintf("note-user-%d-keys", user_id)
	var keys []string
	if cacheKey, found := s.cache.Get(userCacheKey); found {
		keys = cacheKey.([]string)
	}
	keys = append(keys, key)
	s.cache.Set(userCacheKey, keys, cache.DefaultExpiration)
}

func (s *NoteService) InvalidataUserCache(user_id int) {
	userCacheKey := fmt.Sprintf("note-user-%d-keys", user_id)
	if keys, found := s.cache.Get(userCacheKey); found {
		for _, key := range keys.([]string) {
			s.cache.Delete(key)
		}
		s.cache.Set(userCacheKey, []string{}, cache.DefaultExpiration)
	}
}

func validateNoteData(note models.Note) error {
	if note.Title == "" || note.Content == "" {
		return errors.New("title and content must be provided")
	}
	return nil
}

func (s *NoteService) CreateNote(ctx context.Context, note models.Note) error {
	if err := validateNoteData(note); err != nil {
		return err
	}

	noteId, err := s.repo.CreateNote(ctx, note)
	if err != nil {
		return err
	}

	userID, err := s.repo.GetUserByFinance(noteId)
	if err != nil {
		return err
	}
	s.InvalidataUserCache(userID)
	return nil
}

func (s *NoteService) DeleteNote(id_to_del int) error {
	userID, err := s.repo.GetUserByFinance(id_to_del)
	if err != nil {
		return err
	}

	err = s.repo.DeleteNote(context.Background(), id_to_del)
	if err != nil {
		return err
	}

	s.InvalidataUserCache(userID)
	return nil
}

func (s *NoteService) FetchNote(user_id, start, end int, folder_id ...int) ([]models.Note, error) {
	var key string
	if len(folder_id) > 0 {
		key = fmt.Sprintf("note-%d-%d-%d-%d", user_id, start, end, folder_id[0])
	} else {
		key = fmt.Sprintf("note-%d-%d-%d", user_id, start, end)
	}

	if cachedData, found := s.cache.Get(key); found {
		return cachedData.([]models.Note), nil
	}

	data, err := s.repo.FetchNote(context.Background(), user_id, start, end, folder_id...)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, data, cache.DefaultExpiration)
	s.AddCacheKeyForUser(user_id, key)
	return data, nil
}

func (s *NoteService) UpdateNote(ctx context.Context, note models.Note) error {
	userId, err := s.repo.GetUserByFinance(note.Id)
	if err != nil {
		return err
	}

	err = s.repo.UpdateNote(ctx, note)
	if err != nil {
		return err
	}

	s.InvalidataUserCache(userId)
	return nil
}
