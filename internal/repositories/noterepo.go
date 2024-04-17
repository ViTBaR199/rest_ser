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
	FetchNote(ctx context.Context, start, end int, folder_id ...int) ([]models.Note, error)
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

func (r *noteRepositories) FetchNote(ctx context.Context, start, end int, folder_id ...int) ([]models.Note, error) {
	var result []models.Note
	var rows *sql.Rows
	var err error

	if len(folder_id) > 0 {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_notes($1, $2, $3)", start, end, folder_id[0])
	} else {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_notes($1, $2)", start, end)
	}

	if err != nil {
		return nil, fmt.Errorf("querying fetch_notes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.Id, &n.Title, &n.Content, &n.Folder_id, &n.Folder_name); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		result = append(result, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return result, nil
}
