package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/models"

	_ "github.com/lib/pq"
)

type FolderRepositories interface {
	CreateFolder(ctx context.Context, folder models.Folder) error
	DeleteFolder(ctx context.Context, id_to_del int) error
	FetchFolder(ctx context.Context, folder_type string) ([]models.Folder, error)
	UpdateFolder(ctx context.Context, folder models.Folder) error
}

type folderRepositories struct {
	db *sql.DB
}

func NewFolderRepositories(db *sql.DB) FolderRepositories {
	return &folderRepositories{db: db}
}

func (r *folderRepositories) CreateFolder(ctx context.Context, folder models.Folder) error {
	_, err := r.db.ExecContext(ctx, "SELECT create_new_folder($1, $2, $3, $4)", folder.Name, folder.Type, folder.Image, folder.Color)
	return err
}

func (r *folderRepositories) DeleteFolder(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_folder($1)", id_to_del)
	return err
}

func (r *folderRepositories) FetchFolder(ctx context.Context, folder_type string) ([]models.Folder, error) {
	var folders []models.Folder
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM fetch_folders($1)", folder_type)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_folders: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Folder
		var color sql.NullInt64

		if err := rows.Scan(&f.ID, &f.Name, &f.Type, &f.Image, &color, &f.Folder_count); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		if color.Valid {
			tempColor := int(color.Int64)
			f.Color = &tempColor
		} else {
			f.Color = nil
		}

		folders = append(folders, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return folders, nil
}

func (r *folderRepositories) UpdateFolder(ctx context.Context, folder models.Folder) error {
	_, err := r.db.ExecContext(ctx, "SELECT update_folder($1, $2, $3, $4)", folder.ID, folder.Name, folder.Type, folder.Image)
	return err
}
