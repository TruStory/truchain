package graphql

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/samsarahq/thunder/batch"
	thunder "github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/introspection"
	builder "github.com/samsarahq/thunder/graphql/schemabuilder"
	"github.com/samsarahq/thunder/reactive"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

// Request represents the JSON body of a GraphQL query request
type Request struct {
	Query     string                 `json:"query"`     // The GraphQL query string
	Variables map[string]interface{} `json:"variables"` // Variable values for the query
}

// Client holds a GraphQL schema / execution context
type Client struct {
	pendingSchema *builder.Schema
	queries       *builder.Object
	mutations     *builder.Object
	Schema        *thunder.Schema
	Built         bool
}

// NewGraphQLClient returns a GraphQL client with an empty, unbuilt schema
func NewGraphQLClient() *Client {
	schema := builder.NewSchema()
	client := Client{pendingSchema: schema, queries: schema.Query(), mutations: schema.Mutation(), Schema: nil, Built: false}
	return &client
}

// Query runs the GraphQL query in a given `Request` and returns a `chttp.Response` containing a `Response` object
func (c *Client) Query(withCtx context.Context, r Request) chttp.Response {
	if !c.Built {
		c.BuildSchema()
	}
	query, err := c.prepareQuery(r.Query, r.Variables)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	return c.runQuery(withCtx, query)
}

// RegisterQueryResolver adds a top-level resolver to find the first batch of entities in a GraphQL query
func (c *Client) RegisterQueryResolver(name string, fn interface{}) {
	c.queries.FieldFunc(name, fn)
}

// RegisterPaginatedQueryResolver adds a top-level resolver to find the first paginated batch of entities in a GraphQL query
func (c *Client) RegisterPaginatedQueryResolver(name string, fn interface{}) {
	c.queries.FieldFunc(name, fn, builder.Paginated)
}

// RegisterMutation registers a mutation
func (c *Client) RegisterMutation(name string, fn interface{}) {
	c.mutations.FieldFunc(name, fn)
}

// RegisterObjectResolver adds a set of field resolvers for objects of the given type that are returned by top-level resolvers
func (c *Client) RegisterObjectResolver(name string, objPrototype interface{}, fields map[string]interface{}) {
	obj := c.pendingSchema.Object(name, objPrototype)
	for fieldName, fn := range fields {
		obj.FieldFunc(fieldName, fn)
	}
}

// RegisterPaginatedObjectResolver adds a set of paginated field resolvers for objects of the given type that are returned by top-level resolvers
func (c *Client) RegisterPaginatedObjectResolver(name, key string, objPrototype interface{}, fields map[string]interface{}) {
	obj := c.pendingSchema.Object(name, objPrototype)
	obj.Key(key)

	for fieldName, fn := range fields {
		obj.FieldFunc(fieldName, fn)
	}
}

// BuildSchema builds the GraphQL schema from the given resolvers and
func (c *Client) BuildSchema() {
	builtSchema := c.pendingSchema.MustBuild()
	introspection.AddIntrospectionToSchema(builtSchema)
	c.Schema = builtSchema
	c.Built = true
}

// GenerateSchema writes the GraphQL schema to a file
func (c *Client) GenerateSchema() {
	valueJSON, err := introspection.ComputeSchemaJSON(*c.pendingSchema)
	if err != nil {
		panic(err)
	}

	rootdir := viper.GetString(cli.HomeFlag)
	if rootdir == "" {
		rootdir = os.ExpandEnv("$HOME/.truchaind")
	}

	path := filepath.Join(rootdir, "graphql-schema.json")
	err = ioutil.WriteFile(path, valueJSON, 0644)
	if err != nil {
		panic(err)
	}
}

func (c *Client) runQuery(withCtx context.Context, query *thunder.Query) chttp.JSONResponse {
	var wg sync.WaitGroup
	var response chttp.JSONResponse
	e := thunder.Executor{}

	wg.Add(1)

	runner := reactive.NewRerunner(withCtx, func(ctx context.Context) (interface{}, error) {
		defer wg.Done()
		data, err := e.Execute(batch.WithBatching(ctx), c.Schema.Query, nil, query)
		if err != nil {
			response = chttp.SimpleErrorResponse(400, err).(chttp.JSONResponse)
			return nil, err
		}

		rawResBytes, err := json.Marshal(data)
		if err != nil {
			response = chttp.SimpleErrorResponse(500, err).(chttp.JSONResponse)
			return nil, err
		}

		resBytes := []byte(strings.Replace(string(rawResBytes), "iD", "id", -1))
		response = chttp.SimpleResponse(200, resBytes).(chttp.JSONResponse)

		return data, nil
	}, thunder.DefaultMinRerunInterval)

	wg.Wait()
	runner.Stop()

	return response
}

func (c *Client) prepareQuery(qs string, params map[string]interface{}) (*thunder.Query, error) {
	query, err := thunder.Parse(qs, params)
	if err != nil {
		return nil, err
	}

	err = thunder.PrepareQuery(c.Schema.Query, query.SelectionSet)
	if err != nil {
		return nil, err
	}

	return query, nil
}
