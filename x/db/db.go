package db

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

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

// GenericMutations write to the database
type GenericMutations interface {
	Add(model ...interface{}) error
	Update(model interface{}) error
	RegisterModel(model interface{}) error
	Remove(model interface{}) error
}

// Add adds any number of models as a database rows
func (c *Client) Add(model ...interface{}) error {
	return c.Insert(model...)
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

// Update updates a model
func (c *Client) Update(model interface{}) error {
	return c.Update(model)
}

// Remove deletes a models from a table
func (c *Client) Remove(model interface{}) error {
	return c.Delete(model)
}

// GenericQueries are generic reads for models
type GenericQueries interface {
	Count(model interface{}) (int, error)
	Find(model interface{}) error
	FindAll(models interface{}) error
}

// Count returns the count of the model
func (c *Client) Count(model interface{}) (count int, err error) {
	count, err = c.Model(model).Count()

	return
}

// Find selects a single model by primary key
func (c *Client) Find(model interface{}) error {
	return c.Select(model)
}

// FindAll selects all models
func (c *Client) FindAll(models interface{}) error {
	return c.Model(models).Select()
}
