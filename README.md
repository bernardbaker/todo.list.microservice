# Module Overview

Create a gRPC service that allows basic CRUD operations for a Todo list. This will include:

1. Adding a Todo Item: With fields like `id`, `title`, and `completed` (boolean).
2. Listing All Todo Items: Returns all items in the Todo list.
3. Marking a Todo as Completed: Update the `completed` field of a specific item.
4. Deleting a Todo Item: Remove a Todo by its `id`.

# Steps to Implement

1. Define the gRPC Service:

   - Write a `.proto` file with the service definition, including the RPC methods:

```proto
syntax = "proto3";

package todo;

service TodoService {
    rpc AddTodo (TodoItem) returns (TodoResponse);
    rpc ListTodos (Empty) returns (TodoList);
    rpc MarkCompleted (TodoRequest) returns (TodoResponse);
    rpc DeleteTodo (TodoRequest) returns (TodoResponse);
}

message TodoItem {
    string id = 1;
    string title = 2;
    bool completed = 3;
}

message TodoRequest {
    string id = 1;
}

message TodoResponse {
    string message = 1;
}

message TodoList {
    repeated TodoItem todos = 1;
}

message Empty {}
```

2. Generate gRPC Code:

   - Use the `protoc` tool with the Go gRPC plugin to generate the Go code for the service and messages.

3. Implement the Server:

   - Create a `main.go` file and implement the server in Golang.
   - Use an in-memory store (like a slice or a map) for storing the Todo items.
   - Example of a handler:

```go
var todos = make(map[string]*todo.TodoItem)

func (s *server) AddTodo(ctx context.Context, item *todo.TodoItem) (*todo.TodoResponse, error) {
    todos[item.Id] = item
    return &todo.TodoResponse{Message: "Todo added!"}, nil
}
```

4. Implement the Client (Optional):

   - Write a small client in Go to interact with the server for testing.
   - Example of adding a Todo item from the client:

```go
conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
client := todo.NewTodoServiceClient(conn)

item := &todo.TodoItem{Id: "1", Title: "Learn gRPC", Completed: false}
response, _ := client.AddTodo(context.Background(), item)
fmt.Println(response.Message)
```

5. Run and Test:

   - Start the server.
   - Use the client (or `grpcurl` command-line tool) to test the service.
