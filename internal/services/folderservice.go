package services

import (
	"context"
	"errors"
	"myapp/internal/models"
	"myapp/internal/repositories"
)

type FolderService struct {
	repo repositories.FolderRepositories
}

func NewFolderService(repo repositories.FolderRepositories) *FolderService {
	return &FolderService{repo: repo}
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
		return nil
	}
	return s.repo.CreateFolder(ctx, folder)
}

func (s *FolderService) DeleteFolder(id_to_del int) error {
	return s.repo.DeleteFolder(context.Background(), id_to_del)
}

func (s *FolderService) FetchFolder(folder_type string) ([]models.Folder, error) {
	return s.repo.FetchFolder(context.Background(), folder_type)
}

func (s *FolderService) UpdateFolder(ctx context.Context, folder models.Folder) error {
	return s.repo.UpdateFolder(ctx, folder)
}
