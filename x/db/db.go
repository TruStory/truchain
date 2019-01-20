package db

import (
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/joho/godotenv"
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
	// db := pg.Connect(&pg.Options{
	// 	User:     "blockshane",
	// 	Password: "",
	// 	Database: "trudb",
	// })

	// db := pg.Connect(&pg.Options{
	// 	Addr:     "localhost:5432",
	// 	User:     "blockshane",
	// 	Password: "",
	// 	Database: "trudb",
	// })

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

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
