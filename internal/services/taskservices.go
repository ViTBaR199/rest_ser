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
		return nil
	}

	return s.repo.CreateTask(ctx, task)
}

func (s *TaskServices) DeleteTask(id_to_del int) error {
	return s.repo.DeleteTask(context.Background(), id_to_del)
}

func (s *TaskServices) FetchTask(start, end int) ([][]string, error) {
	return s.repo.FetchTask(context.Background(), start, end)
}
