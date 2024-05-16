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
	FetchFolder(ctx context.Context, start, end int, id_user int, type_folder string) ([]models.Folder, error)
	UpdateFolder(ctx context.Context, folder models.Folder) error
	GetUserByFolder(folderID int) (int, error)
	FetchFolderById(ctx context.Context, id_folder, user_id int) ([]models.Folder, error)
}

type folderRepositories struct {
	db *sql.DB
}

func NewFolderRepositories(db *sql.DB) FolderRepositories {
	return &folderRepositories{db: db}
}

func (r *folderRepositories) CreateFolder(ctx context.Context, folder models.Folder) error {
	_, err := r.db.ExecContext(ctx, "SELECT create_new_folder($1, $2, $3, $4, $5)", folder.Name, folder.Type, folder.Image, folder.Color, folder.User_id)
	return err
}

func (r *folderRepositories) DeleteFolder(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_folder($1)", id_to_del)
	return err
}

func (r *folderRepositories) FetchFolder(ctx context.Context, start, end int, id_user int, type_folder string) ([]models.Folder, error) {
	var folders []models.Folder
	var rows *sql.Rows
	var err error

	if len(type_folder) > 0 {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_folders($1, $2, $3, $4)", start, end, type_folder, id_user)
	} else {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_folders($1, $2, $3)", start, end, id_user)
	}

	if err != nil {
		return nil, fmt.Errorf("querying fetch_task: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Folder
		var color sql.NullInt64

		if err := rows.Scan(&f.ID, &f.Name, &f.Type, &f.Image, &color, &f.Folder_count, &f.User_id); err != nil {
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
	if folder.ID != nil && folder.Name != "" && folder.Type != "" &&
		folder.Image != "" && folder.Color != nil && folder.User_id != nil {
		_, err := r.db.ExecContext(ctx, "SELECT update_folder($1, $2, $3, $4, $5, $6)", folder.ID, folder.Name, folder.Type, folder.Image, folder.Color, folder.User_id)
		return err
	} else if folder.ID == nil {
		return fmt.Errorf("the ID are set incorrectly")
	} else if folder.Name == "" {
		return fmt.Errorf("the NAME are set incorrectly")
	} else if folder.Type == "" {
		return fmt.Errorf("the TYPE are set incorrectly")
	} else if folder.Image == "" {
		return fmt.Errorf("the IMAGE are set incorrectly")
	} else if folder.Color == nil {
		return fmt.Errorf("the COLOR are set incorrectly")
	} else if folder.User_id == nil {
		return fmt.Errorf("the USER ID are set incorrectly")
	}
	return fmt.Errorf("the parameters are set incorrectly")
}

func (r *folderRepositories) GetUserByFolder(folderID int) (int, error) {
	var userId int
	err := r.db.QueryRow("SELECT get_user_by_folderID($1)", folderID).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (r *folderRepositories) FetchFolderById(ctx context.Context, id_folder, user_id int) ([]models.Folder, error) {
	var folders []models.Folder
	var rows *sql.Rows
	var err error
	rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_folder_by_id($1, $2)", id_folder, user_id)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_task: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Folder
		var color sql.NullInt64

		if err := rows.Scan(&f.ID, &f.Name, &f.Type, &f.Image, &color, &f.Folder_count, &f.User_id); err != nil {
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
