package db

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Datastore defines all operations on the DB
// This interface can be mocked out for tests, etc.
type Datastore interface {
	Mutations
	Queries
}

// Mutations write to the database
type Mutations interface {
	Add(model interface{}) error
	RegisterModel(model interface{}) error
}

// Client is a Postgres client.
// It wraps a pool of Postgres DB connections.
type Client struct {
	*pg.DB
}

// NewDBClient creates a Postgres client
func NewDBClient() *Client {
	db := pg.Connect(&pg.Options{
		Addr:     os.Getenv("PG_ADDR"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_USER_PW"),
		Database: os.Getenv("PG_DB_NAME"),
	})

	return &Client{db}
}

// Add implements `Datastore`.
// It adds a model as a database row.
func (c *Client) Add(model interface{}) error {
	return c.Insert(model)
}

// RegisterModel creates a table for a type.
// A table is automatically created based on the passed in struct fields.
func (c *Client) RegisterModel(model interface{}) error {
	fmt.Printf("Registered model %v\n", model)
	return c.CreateTable(model, &orm.CreateTableOptions{
		Temp:        false,
		IfNotExists: true,
	})
}
