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
	FetchFolder(ctx context.Context, start, end int) ([][]string, error)
	UpdateFolder(ctx context.Context, folder models.Folder) error
}

type folderRepositories struct {
	db *sql.DB
}

func NewFolderRepositories(db *sql.DB) FolderRepositories {
	return &folderRepositories{db: db}
}

func (r *folderRepositories) CreateFolder(ctx context.Context, folder models.Folder) error {
	_, err := r.db.ExecContext(ctx, "SELECT create_new_folder($1, $2, $3)", folder.Name, folder.Type, folder.Image)
	return err
}

func (r *folderRepositories) DeleteFolder(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_folder($1)", id_to_del)
	return err
}

func (r *folderRepositories) FetchFolder(ctx context.Context, start, end int) ([][]string, error) {
	var result [][]string

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM fetch_folders($1, $2)", start, end)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_folders: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, folder_type, image string
		if err := rows.Scan(&id, &name, &folder_type, &image); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		rowData := []string{fmt.Sprintf("%d", id), name, folder_type, image}
		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return result, nil
}

func (r *folderRepositories) UpdateFolder(ctx context.Context, folder models.Folder) error {
	_, err := r.db.ExecContext(ctx, "SELECT update_folder($1, $2, $3, $4)", folder.ID, folder.Name, folder.Type, folder.Image)
	return err
}
