# TruAPI

This file contains step-by-step instructions for extending the server in a few of the most common ways; namely, supporting new messages, queries, and GraphQL resolvers.

## To support a new message

- Ensure message type exists on server
- Ensure there is a route to handle messages from that package
- Enable the API to decode this message type
  - In `x/truapi/supported_msgs.go`
      - Add an entry in the `supported` map for your new message type. The key is the name of the type, and the value is an empty struct of that type.
- Enable the chain to handle this message type
  - In `x/$msgpkg/handler.go`
    - Add function to handle message type
    - If necessary, modify the `NewHandler` call site to pass additional keepers
    - Add case for new message type inside `NewHandler`

## To support a new querier/query

Queriers are string-dispatched functions that are configured with keepers to access the blockchain (see example in `x/category/queryable.go`). Queries are how non-blockchain constructs (like GraphQL resolvers called from HTTP handlers) access blockchain state.

### To add a Querier to a package
  - Create the file `mypackage/queryable.go` and add the Querier code here
  - Organize the code like you see in other `queryable.go` files:
    - The `QueryPath` constant describes the root string the Querier will be dispatched by
    - The `NewQuerier` method takes the necessary keepers, and returns a function matching the `sdk.Querier` signature
      - The body of this function is a big string switch, dispatching from `path[0]`
  - Mount the Querier in `app/app.go` like so:

    ```go
    app.QueryRouter()
      .AddRoute(story.QueryPath, story.NewQuerier(app.storyKeeper))
    ```

### To add a new query to an existing Querier
  - Add a constant holding the public path element of this query (e.g. `StoriesByCategoryID = "category"`)
    - The name of this constant should match the name of the function implementing the query
  - Add a function that implements this query
  - Add a case to the Querier's switch, which will call this function with necessary keepers and request data
  - *Important note:* Because Go-Amino represents numbers as strings, numbers must be strings in the query parameters and parsed in the Querier before being passed on as Keeper method args (see `x/category/queryable.go`).


## To support a new GraphQL resolver

GraphQL resolvers can be top-level queries (e.g. `stories(categoryID)`), or struct fields (e.g. `category { stories { age } }`)

- Identify the source of the data you need to resolve
  - If it's an external DB or API, make sure there are configs for those connections
  - If it's the blockchain, make sure there is (and if not, create) a Querier with access to that data
    - To call a Querier in `resolver.go`, call `ta.RunQuery(mypackage.QueryRoot + "/" + mypackage.MyQueryPath)`

- Identify the parameters
  - For top-level queries, these will be the arguments like `categoryID` above
  - For struct fields, these will be an instance of the struct on which the field is being defined

- Write a resolver
  - See `x/truapi/resolver.go` for examples of both resolver types (and add your resolver to this file)
  - A query resolver has the signature `func(context.Context, mypackage.MyQueryParamStructType)`
  - A struct field resolver has the signature `func(context.Context, mypackage.MyStructType)`

- Register the resolver
  - In `x/truapi/truapi.go`, call the appropriate registration method inside of `RegisterResolvers()`
  - To register a query resolver, simply call `ta.GraphQLClient.RegisterQueryResolver("myQuery", ta.myQueryResolver)`
  - All fields for a given struct are registered together, like so:
    ```go
    ta.GraphQLClient.RegisterObjectResolver("Story", story.Story{}, map[string]interface{}{
      "id": func(_ context.Context, q story.Story) int64 {
        return q.ID
      },
      "category": ta.storyCategoryResolver,
    })
    ```
