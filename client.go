package orm

import "database/sql"

type Client struct {
	DB *sql.DB
}

// func (c *Client) Tx() *TxSession {
// 	return &TxSession{}
// }

func (c *Client) Table(schema Schema) *Action {
	return &Action{
		db:     c.DB,
		schema: schema,
	}
}
