# crdb-grpc
A barebones example of accessing CockroachDB via GraphQL

### Dependencies

* [Go](https://go.dev)

### Running locally

Create a cluster
```
$ cockroach start-single-node \
		--listen-addr=localhost:26257 \
		--http-addr=localhost:8080 \
		--insecure
```

Create a table
```
$ cockroach sql --insecure < create.sql
```

Start the server
```
$ go run main.go
```

### Sample requests

Fetch all todos
``` graphql
{
  todos {
    id
    title
  }
}
```

Fetch one todo
``` graphql
{
  todo(id: "8da8291c-5985-41c9-8069-0de865dd20d7") {
    title
  }
}
```

Create a todo
``` graphql
mutation CreateTodo($todo: TodoInput!) {
  createTodo(todo: $todo) {
    id
    title
  }
}
```

``` json
{
  "todo": {
  	"title": "todo d"
  }
}
```

Delete a todo
``` graphql
mutation DeleteTodo($id: ID!) {
  deleteTodo(id: $id)
}
```

```
{
  "id": "8da8291c-5985-41c9-8069-0de865dd20d7"
}
```