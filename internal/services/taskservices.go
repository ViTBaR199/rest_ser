package services

import (
	"context"
	"errors"
	"myapp/internal/models"
	"myapp/internal/repositories"
)

type TaskServices struct {
	repo repositories.TaskRepositories
}

func NewTaskServices(repo repositories.TaskRepositories) *TaskServices {
	return &TaskServices{repo: repo}
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
		//-----------изменить возвращаемое значение--------------
		return err
	}

	return s.repo.CreateTask(ctx, task)
}

func (s *TaskServices) DeleteTask(id_to_del int) error {
	return s.repo.DeleteTask(context.Background(), id_to_del)
}

func (s *TaskServices) FetchTask(start, end int, folder_id ...int) ([]models.Task, error) {
	return s.repo.FetchTask(context.Background(), start, end, folder_id...)
}

func (s *TaskServices) UpdateTask(ctx context.Context, task models.Task) error {
	return s.repo.UpdateTask(ctx, task)
}

func (s *TaskServices) CountTask() (int, error) {
	return s.repo.CountTask(context.Background())
}

func (s *TaskServices) CountTaskFavourites() (int, error) {
	return s.repo.CountTaskFavourites(context.Background())
}

func (s *TaskServices) FetchTaskFavourites(start, end int, folder_id ...int) ([]models.Task, error) {
	return s.repo.FetchTaskFavourites(context.Background(), start, end, folder_id...)
}
