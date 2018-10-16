package graphql

import (
  "context"
  "encoding/json"
  "sync"

  "github.com/samsarahq/thunder/batch"
  "github.com/samsarahq/thunder/graphql/introspection"
  "github.com/samsarahq/thunder/reactive"
  "github.com/TruStory/truchain/x/chttp"

  builder "github.com/samsarahq/thunder/graphql/schemabuilder"
  thunder "github.com/samsarahq/thunder/graphql"
)

type Request struct {
  Query string                     `json:"query"`    // The GraphQL query string
  Variables map[string]interface{} `json:"variables` // Variable values for the query
}

type Response struct {
  Data interface{} `json:"data"`
}

type Client struct {
  pendingSchema *builder.Schema
  queries *builder.Object
  Schema *thunder.Schema
  Built bool
}

func NewGraphqlClient() *Client {
  schema := builder.NewSchema()
  client := Client{pendingSchema: schema, queries: schema.Query(), Schema: nil, Built: false}
  return &client
}

func (c *Client) Query(withCtx context.Context, r Request) chttp.Response {
  if !c.Built {
    c.BuildSchema()
  }

  query, err := c.prepareQuery(r.Query, r.Variables)

  if err != nil {
    return chttp.SimpleErrorResponse(401, err)
  }

  return c.runQuery(withCtx, query)
}

func (c *Client) RegisterQueryResolver(name string, fn interface{}) {
  c.queries.FieldFunc(name, fn)
}

func (c *Client) RegisterObjectResolver(name string, objPrototype interface{}, fields map[string]interface{}) {
  obj := c.pendingSchema.Object(name, objPrototype)

  for fieldName, fn := range fields {
    obj.FieldFunc(fieldName, fn)
  }
}

func (c *Client) BuildSchema() {
  builtSchema := c.pendingSchema.MustBuild()
  introspection.AddIntrospectionToSchema(builtSchema)
  c.Schema = builtSchema
  c.Built = true
}

func (c *Client) runQuery(withCtx context.Context, query *thunder.Query) chttp.JsonResponse {
  var wg sync.WaitGroup
  var response chttp.JsonResponse
  e := thunder.Executor{}

  wg.Add(1)

  runner := reactive.NewRerunner(withCtx, func(ctx context.Context) (interface{}, error) {
    defer wg.Done()
    data, err := e.Execute(batch.WithBatching(ctx), c.Schema.Query, nil, query)
    
    if err != nil {
      response = chttp.SimpleErrorResponse(400, err).(chttp.JsonResponse)
      return nil, err
    }

    resBytes, err := json.Marshal(data)

    if err != nil {
      response = chttp.SimpleErrorResponse(500, err).(chttp.JsonResponse)
      return nil, err
    }

    response = chttp.SimpleResponse(200, resBytes).(chttp.JsonResponse)

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

  if err != nil {
    return nil, err
  }

  err = thunder.PrepareQuery(c.Schema.Query, query.SelectionSet)

  if err != nil {
    return nil, err
  }

  return query, nil
}
