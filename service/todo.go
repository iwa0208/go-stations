package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// In normal use, create one Stmt when your process starts. 確認作業？
	// stmt, err := s.db.PrepareContext(ctx, insert)
	// if err != nil {
	// 	// log.Printf(err.Error())
	// 	return nil, err
	// }
	// defer stmt.Close()

	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		// log.Printf(err.Error())
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		// log.Printf(err.Error())
		return nil, err
	}

	var created, updated time.Time
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&subject, &description, &created, &updated) // rowの形式で返ってくるものを 構造体todoに書き換えて返却する
	if err != nil {
		// log.Printf(err.Error())
		return nil, err
	}
	return &model.TODO{
		ID:          int64(id),
		Subject:     subject,
		Description: description,
		CreatedAt:   created,
		UpdatedAt:   updated,
	}, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var ttodo []*model.TODO //空行列

	// sizeは存在しない場合は5に設定する。predIDが存在しないときはないものとみなす。
	if size == 0 {
		size = 5
	}
	if prevID == 0 {
		rows, err := s.db.QueryContext(ctx, read, size)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var subject, description string
			var created, updated time.Time
			if err := rows.Scan(&id, &subject, &description, &created, &updated); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				log.Fatal(err)
			}
			ttodo = append(ttodo, &model.TODO{
				ID:          int64(id),
				Subject:     subject,
				Description: description,
				CreatedAt:   created,
				UpdatedAt:   updated,
			})
		}
	} else {
		rows, err := s.db.QueryContext(ctx, readWithID, prevID, size)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var subject, description string
			var created, updated time.Time
			if err := rows.Scan(&id, &subject, &description, &created, &updated); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				log.Fatal(err)
			}
			ttodo = append(ttodo, &model.TODO{
				ID:          int64(id),
				Subject:     subject,
				Description: description,
				CreatedAt:   created,
				UpdatedAt:   updated,
			})
		}
	}

	return ttodo, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		// log.Fatalf("expected to affect 1 row, affected %d", rows)
		err = &model.ErrNotFound{
			What: "対象のレコードはありませんでした。",
		}
		return nil, err
	}

	var created, updated time.Time
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&subject, &description, &created, &updated) // rowの形式で返ってくるものを 構造体todoに書き換えて返却する
	if err != nil {
		return nil, err
	}
	return &model.TODO{
		ID:          int64(id),
		Subject:     subject,
		Description: description,
		CreatedAt:   created,
		UpdatedAt:   updated,
	}, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?)`

	for _, id := range ids {
		_, err := s.db.ExecContext(ctx, deleteFmt, id)
		if err != nil {
			return err
		}
	}
	return nil
}
