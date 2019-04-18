# TruChain: DB Module

This module wraps a Postgres database.

The database is accessed via an interface that defines reads and writes:

```go
type Datastore interface {
	Mutations
	Queries
}
```

## Installation and Setup

On macOS:

```sh
# install postgres
brew update
brew install postgresql

# start postgres on launch
brew services start postgresql

# create a test database
createdb trudb
```

With Docker:

```sh
# Starts a new container in the background and creates `trudb` database
docker run --name trudb \
-e POSTGRES_DB=trudb \
-e POSTGRES_USER=postgres \
-e POSTGRES_PASSWORD=postgres \
-p 5432:5432 \
-d postgres:11.1
```

Create a `.env` with the following variables for local setup:

```
PG_ADDR=localhost:5432
PG_USER=[your local machine user from `whoami`]
PG_USER_PW=
PG_DB_NAME=trudb
```

Copy `.env` to the default chain config path, `.chain` locally, and `~/.truchaind` on AWS.

## Mutations

Writes conform to the interface:

```go
type Mutations interface {
	Add(model interface{}) error
}
```

### Add a row

Call `Add(model interface{})` with the data type you want to persist. An auto-incrementing `ID` primary key will automatically be created.

## Queries

Add queries to the `Queries` interface:

```go
type Queries interface {
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

## Notes

Refer to [https://github.com/go-pg/pg](https://github.com/go-pg/pg) for advanced features.
