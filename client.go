package orm

import (
	"context"
	"database/sql"
	"fmt"
)

type Client struct {
	DB         *sql.DB
	driverName string
}

func (c *Client) Transaction(ctx context.Context, fn func(s *Session) error) (err error) {
	tx, err := c.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				err = fmt.Errorf("rollback error: %w, original error: %w", txErr, err)
			}
		} else {
			err = tx.Commit()
		}
	}()
	s := &Session{db: tx}
	return fn(s)
}

func (c *Client) Table(schema Schema) *Action {
	s := &Session{db: c.DB}
	return s.Table(schema)
}

func NewClient(driverName, dataSourceName string) (*Client, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Client{DB: db, driverName: driverName}, nil
}
