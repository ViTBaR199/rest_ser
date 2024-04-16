package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/models"

	_ "github.com/lib/pq"
)

type NoteRepositories interface {
	CreateNote(ctx context.Context, note models.Note) error
	DeleteNote(ctx context.Context, id_to_del int) error
	FetchNote(ctx context.Context, start, end int) ([][]string, error)
}

type noteRepositories struct {
	db *sql.DB
}

func NewNoteRepositories(db *sql.DB) NoteRepositories {
	return &noteRepositories{db: db}
}

func (r *noteRepositories) CreateNote(ctx context.Context, note models.Note) error {
	_, err := r.db.ExecContext(ctx, "SELECT FROM create_new_note($1, $2, $3)", note.Title, note.Content, note.Folder_id)
	return err
}

func (r *noteRepositories) DeleteNote(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_note($1)", id_to_del)
	return err
}

func (r *noteRepositories) FetchNote(ctx context.Context, start, end int) ([][]string, error) {
	var result [][]string

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM fetch_notes($1, $2)", start, end)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_notes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, folder_id int
		var title, content string
		if err := rows.Scan(&id, &title, &content, &folder_id); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		rowData := []string{fmt.Sprintf("%d", id), title, content, fmt.Sprintf("%d", folder_id)}
		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return result, nil
}
