# TruChain: DB Module

This module wraps a Postgres database.

### Create a new table

Call `RegisterModel(model interface{})` in `TruAPI.RegisterModels()` where `model` is a struct with all the fields for the table. A table will automatically
be created based on the struct fields. Go fields will be automatically translated to
Postgres data types.

### Add a row

Call `Add(model interface{})` with the data type you want to persist. An auto-incrementing `ID` primary key will automatically be created.

### Queries

Add queries to the `Datastore` interface:

```go
type Datastore interface {
	Add(model interface{}) error
	RegisterModel(model interface{}) error
	TwitterProfileByAddress(addr string) (TwitterProfile, error)
}
```

This interface can be mocked out for testing.

### GraphQL

The Postgres client has been added to `TruAPI`, which means it can be accessed in GraphQL resolvers. 

For example, the `TwitterProfile` type is resolved with: 

```go
	ta.GraphQLClient.RegisterObjectResolver("TwitterProfile", db.TwitterProfile{}, map[string]interface{}{
		"id": func(_ context.Context, q db.TwitterProfile) string { return string(q.ID) },
	})
```

The GraphQL client is agnostic of data source. Resolvers can access data on the chain and/or data in Postgres, without the mobile or web client knowing anything about the underlying data source.

That's all folks! Easy peasy.

### Notes

Database tables are only created in the genesis block. If you add new tables, the chain has to be restarted, or a migration needs to be created that registers models at a later time.

Refer to [https://github.com/go-pg/pg](https://github.com/go-pg/pg) for advanced features.
