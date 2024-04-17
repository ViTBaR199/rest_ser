package services

import (
	"context"
	"errors"
	"myapp/internal/models"
	"myapp/internal/repositories"
)

type NoteService struct {
	repo repositories.NoteRepositories
}

func NewNoteService(repo repositories.NoteRepositories) *NoteService {
	return &NoteService{repo: repo}
}

func validateNoteData(note models.Note) error {
	if note.Title == "" || note.Content == "" {
		return errors.New("title and content must be provided")
	}
	return nil
}

func (s *NoteService) CreateNote(ctx context.Context, note models.Note) error {
	if err := validateNoteData(note); err != nil {
		return nil
	}
	return s.repo.CreateNote(ctx, note)
}

func (s *NoteService) DeleteNote(id_to_del int) error {
	return s.repo.DeleteNote(context.Background(), id_to_del)
}

func (s *NoteService) FetchNote(start, end int, folder_id ...int) ([]models.Note, error) {
	return s.repo.FetchNote(context.Background(), start, end, folder_id...)
}
