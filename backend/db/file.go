package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

var (
	saveFileQuery      = `INSERT INTO files (id, user_id, file_name, content_type, content) VALUES (?, ?, ?, ?, ?)`
	getFileQuery       = `SELECT id, user_id, file_name, content_type, content, created_at, updated_at FROM files WHERE id = ? AND user_id = ?`
	getFilesQuery      = `SELECT id, user_id, file_name, content_type, created_at, updated_at FROM files WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	getFilesCountQuery = `SELECT COUNT(*) FROM files WHERE user_id = ?`
	deleteFileQuery    = `DELETE FROM files WHERE id = ? AND user_id = ?`
	updateFileQuery    = `UPDATE files SET file_name = ?, content_type = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?`
)

type File struct {
	ID          string    `json:"id"`
	UserID      string    `json:"-"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	Content     string    `json:"content,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var (
	ErrFileNotFound = errors.New("file not found")
)

func (r *sqliteRepo) SaveFile(ctx context.Context, userId, fileName, contentType, content string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info("Saving file", "fileName", fileName, "userId", userId)
	id := uuid.New().String()

	_, err := r.Db.ExecContext(ctx, saveFileQuery, id, userId, fileName, contentType, content)
	if err != nil {
		r.logger.Error("failed to save file", "fileName", fileName, "userId", userId, "error", err)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return id, nil
}

// GetFile retrieves a single file by its ID and user ID.
func (r *sqliteRepo) GetFile(ctx context.Context, fileId, userId string) (*File, error) {
	r.logger.Info("Fetching file", "fileId", fileId, "userId", userId)
	row := r.Db.QueryRowContext(ctx, getFileQuery, fileId, userId)

	var file File
	err := row.Scan(&file.ID, &file.UserID, &file.FileName, &file.ContentType, &file.Content, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error("file not found", "fileId", fileId, "userId", userId)
			return nil, ErrFileNotFound
		}
		r.logger.Error("failed to scan file", "fileId", fileId, "userId", userId, "error", err)
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	return &file, nil
}

func (r *sqliteRepo) GetFiles(ctx context.Context, userId string, limit, offset int) ([]*File, error) {
	r.logger.Info("Fetching files for user", "userId", userId, "limit", limit, "offset", offset)
	rows, err := r.Db.QueryContext(ctx, getFilesQuery, userId, limit, offset)

	if err != nil {
		r.logger.Error("failed to query files", "userId", userId, "error", err)
		return nil, fmt.Errorf("failed to query files: %w", err)
	}

	defer rows.Close()

	var files []*File
	for rows.Next() {
		var file File
		err := rows.Scan(&file.ID, &file.UserID, &file.FileName, &file.ContentType, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			r.logger.Error("failed to scan file row", "error", err)
			return nil, fmt.Errorf("failed to scan file row: %w", err)
		}
		files = append(files, &file)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating file rows", "error", err)
		return nil, fmt.Errorf("error iterating file rows: %w", err)
	}

	return files, nil
}

func (r *sqliteRepo) GetFilesCount(ctx context.Context, userId string) (int, error) {
	r.logger.Info("Getting files count for user", "userId", userId)
	row := r.Db.QueryRowContext(ctx, getFilesCountQuery, userId)

	var count int
	err := row.Scan(&count)
	if err != nil {
		r.logger.Error("failed to scan files count", "userId", userId, "error", err)
		return 0, fmt.Errorf("failed to get files count: %w", err)
	}

	return count, nil
}

func (r *sqliteRepo) DeleteFile(ctx context.Context, fileId, userId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info("Deleting file", "fileId", fileId, "userId", userId)
	res, err := r.Db.ExecContext(ctx, deleteFileQuery, fileId, userId)
	if err != nil {
		r.logger.Error("failed to delete file", "fileId", fileId, "userId", userId, "error", err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected after delete", "fileId", fileId, "error", err)
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrFileNotFound
	}

	return nil
}

func (r *sqliteRepo) UpdateFile(ctx context.Context, fileId, userId, fileName, contentType, content string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info("Updating file", "fileId", fileId, "userId", userId)
	res, err := r.Db.ExecContext(ctx, updateFileQuery, fileName, contentType, content, fileId, userId)

	if err != nil {
		r.logger.Error("failed to update file", "fileId", fileId, "error", err)
		return fmt.Errorf("failed to update file: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected after update", "fileId", fileId, "error", err)
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrFileNotFound
	}

	return nil
}
