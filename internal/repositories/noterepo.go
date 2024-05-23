package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/models"

	_ "github.com/lib/pq"
)

type NoteRepositories interface {
	CreateNote(ctx context.Context, note models.Note) (int, error)
	DeleteNote(ctx context.Context, id_to_del int) error
	FetchNote(ctx context.Context, user_id, start, end int, folder_id ...int) ([]models.Note, error)
	UpdateNote(ctx context.Context, note models.Note) error
	GetUserByNote(noteID int) (int, error)
	FetchNoteById(ctx context.Context, noteId int) (models.Note, error)
}

type noteRepositories struct {
	db *sql.DB
}

func NewNoteRepositories(db *sql.DB) NoteRepositories {
	return &noteRepositories{db: db}
}

func (r *noteRepositories) CreateNote(ctx context.Context, note models.Note) (int, error) {
	var newId int
	err := r.db.QueryRowContext(ctx, "SELECT create_new_note($1, $2, $3)", note.Title, note.Content, note.Folder_id).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (r *noteRepositories) DeleteNote(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_note($1)", id_to_del)
	if err != nil {
		return err
	}
	return nil
}

func (r *noteRepositories) FetchNote(ctx context.Context, user_id, start, end int, folder_id ...int) ([]models.Note, error) {
	var result []models.Note
	var rows *sql.Rows
	var err error

	if len(folder_id) > 0 {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_notes($1, $2, $3, $4)", user_id, start, end, folder_id[0])
	} else {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_notes($1, $2, $3)", user_id, start, end)
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

func (r *noteRepositories) GetUserByNote(noteID int) (int, error) {
	var userId int
	err := r.db.QueryRow("SELECT get_user_by_note($1)", noteID).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (r *noteRepositories) UpdateNote(ctx context.Context, note models.Note) error {
	_, err := r.db.ExecContext(ctx, "SELECT update_note($1, $2, $3)", note.Id, note.Title, note.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *noteRepositories) FetchNoteById(ctx context.Context, noteId int) (models.Note, error) {
	var result models.Note
	var rows *sql.Rows
	var err error

	rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_note_by_id($1)", noteId)
	if err != nil {
		return models.Note{}, fmt.Errorf("querying fetch_note: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Note

		if err := rows.Scan(&n.Id, &n.Title, &n.Content, &n.Folder_id); err != nil {
			return models.Note{}, fmt.Errorf("scanning row: %v", err)
		}

		result = n
	}

	return result, nil
}
