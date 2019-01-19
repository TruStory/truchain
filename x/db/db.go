package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Datastore defines all operations on the DB
// This interface can be mocked out for tests, etc.
type Datastore interface {
	Add(model interface{}) error
	RegisterModel(model interface{}) error
	TwitterProfileByAddress(addr string) (TwitterProfile, error)
}

// Client is a Postgres client.
// It wraps a pool of Postgres DB connections.
type Client struct {
	*pg.DB
}

// NewDBClient creates a Postgres client
func NewDBClient() *Client {
	db := pg.Connect(&pg.Options{
		User:     "blockshane",
		Password: "",
		Database: "trudb",
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
	err := c.DropTable(model, &orm.DropTableOptions{
		IfExists: true,
		Cascade:  true,
	})
	if err != nil {
		panic(err)
	}

	return c.CreateTable(model, &orm.CreateTableOptions{
		Temp: false,
	})
}
