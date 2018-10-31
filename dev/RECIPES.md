# Recipes
This file contains step-by-step instructions for extending the server in a few of the most common ways; namely, supporting new messages and queries.

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


## To support a new query

TODO