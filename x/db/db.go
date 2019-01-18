package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Client is a Postgres client
type Client struct {
	*pg.DB
}

// Add adds a model as a database row
func (c *Client) Add(model interface{}) error {
	return c.Insert(model)
}

// Find finds and populates the given model
func (c *Client) Find(model interface{}) error {
	return c.Select(model)
}

// NewDBClient creates a Postgres client
func NewDBClient() *Client {
	db := pg.Connect(&pg.Options{
		User: "blockshane",
	})

	return &Client{db}
}

// RegisterModel creates a table for a type
func (c *Client) RegisterModel(model interface{}) error {
	return c.CreateTable(model, &orm.CreateTableOptions{
		Temp: false,
	})
}
