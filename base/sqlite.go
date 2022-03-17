package base

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	handle *sql.DB
}

type SQLiteTX struct {
	tx *sql.Tx
}

func (db *SQLite) Open(dbName string) error {
	handle, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	db.handle = handle
	return nil
}

func (db *SQLite) Close() error {
	if db.handle != nil {
		return db.handle.Close()
	}

	return nil
}

type ScanHandle func(rows *sql.Rows) error

var (
	InvalidDBHandle = errors.New("invalid handle")
)

func (db *SQLite) Query(handle ScanHandle, query string, args ...interface{}) error {
	if db.handle == nil {
		return InvalidDBHandle
	}

	rows, err := db.handle.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return handle(rows)
}

func (db *SQLite) Exec(query string, args ...interface{}) (int64, int64, error) {
	if db.handle == nil {
		return -1, -1, InvalidDBHandle
	}

	rs, err := db.handle.Exec(query, args...)
	if err != nil {
		return -1, -1, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return -1, -1, err
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		return -1, -1, err
	}

	return id, rows, nil
}

func (db *SQLite) ListTables() ([]string, error) {
	rs := make([]string, 0)
	err := db.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			var table string
			err := rows.Scan(&table)
			if err != nil {
				return err
			}
			rs = append(rs, table)
		}
		return nil
	}, "SELECT tbl_name FROM sqlite_master WHERE type = 'table';")

	return rs, err
}

func (db *SQLite) GetTableRowCount(table string) (int64, error) {
	var total int64
	err := db.Query(func(rows *sql.Rows) error {
		if rows.Next() {
			err := rows.Scan(&total)
			if err != nil {
				return err
			}

			return nil
		}
		return nil
	}, "select count(1) as total from "+table)

	return total, err
}

func (db *SQLite) Begin() (*SQLiteTX, error) {
	if db.handle == nil {
		return nil, InvalidDBHandle
	}

	t, err := db.handle.Begin()
	if err != nil {
		return nil, err
	}
	return &SQLiteTX{tx: t}, nil
}

func (db *SQLiteTX) Commit() error {
	if db.tx == nil {
		return InvalidDBHandle
	}

	err := db.tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (db *SQLiteTX) Rollback() error {
	if db.tx == nil {
		return InvalidDBHandle
	}

	err := db.tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func (db *SQLiteTX) Exec(query string, args ...interface{}) (int64, int64, error) {
	if db.tx == nil {
		return -1, -1, InvalidDBHandle
	}

	rs, err := db.tx.Exec(query, args...)
	if err != nil {
		return -1, -1, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return -1, -1, err
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		return -1, -1, err
	}

	return id, rows, nil
}

func (db *SQLiteTX) Query(handle ScanHandle, query string, args ...interface{}) error {
	if db.tx == nil {
		return InvalidDBHandle
	}

	rows, err := db.tx.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return handle(rows)
}
