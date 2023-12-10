package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
)

type todo struct {
	Id string `json:"id"`
	Todo string `json:"todo"`
	Timestamp int64 `json:"timestamp"`
}
var todos []todo = []todo{}

func responseTodos(w http.ResponseWriter) {
	jsonData, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func todosEndpoint(w http.ResponseWriter, req *http.Request) {
	if(req.Method != "GET") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	responseTodos(w)
}

func addTodoEndpoint(w http.ResponseWriter, req *http.Request) {
	if(req.Method != "POST") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body, err = io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "todo not found", http.StatusBadRequest)
		return
	}
	var newTodo todo

	newTodo.Todo = string(body)
	newTodo.Id = uuid.New().String()
	newTodo.Timestamp = time.Now().Unix()

	todos = append(todos, newTodo)
	
	responseTodos(w)
}

func deleteTodoEndpoint(w http.ResponseWriter, req *http.Request) {
	if(req.Method != "POST") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var todoId, err = io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "todoId not found", http.StatusBadRequest)
		return
	}
	var deleteItemindex int
	for i := 0; i < len(todos); i++ {
		if(todos[i].Id == string(todoId)) {
			deleteItemindex = i
			break
		}
	}
	if deleteItemindex < 0 || deleteItemindex >= len(todos) {
		http.Error(w, "invalid delete index", http.StatusBadRequest)
		return
    }
	todos = append(todos[:deleteItemindex], todos[deleteItemindex+1:]...)

	responseTodos(w)
}

func editTodoEndpoint(w http.ResponseWriter, req *http.Request) {
	if(req.Method != "POST") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get data from req
	var jsonData, err = io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "todoId not found", http.StatusBadRequest)
		return
	}

	// Convert from body to json
	type TodoUpdate struct {
		TodoID  string `json:"todoId"`
		NewValue string `json:"newValue"`
	}

	var todoUpdate TodoUpdate

	err = json.Unmarshal([]byte(jsonData), &todoUpdate)
	if err != nil {
		http.Error(w, "error getting data from request body", http.StatusInternalServerError)
		return
	}

	// update in todos
	var updateItemindex int
	for i := 0; i < len(todos); i++ {
		if(todos[i].Id == string(todoUpdate.TodoID)) {
			updateItemindex  = i
			break
		}
	}
	if updateItemindex < 0 || updateItemindex  >= len(todos) {
		http.Error(w, "invalid update index", http.StatusBadRequest)
		return
    }
	todos[updateItemindex].Todo = todoUpdate.NewValue

	responseTodos(w)
}

func main() {
	const port string = "3001"

	var corsHandler *cors.Cors = cors.Default()
	http.HandleFunc("/todos", func(w http.ResponseWriter, req *http.Request) {
		var handler http.Handler = corsHandler.Handler(http.HandlerFunc(todosEndpoint))
		handler.ServeHTTP(w, req)
	})
	http.HandleFunc("/addTodo", func(w http.ResponseWriter, req *http.Request) {
		var handler http.Handler = corsHandler.Handler(http.HandlerFunc(addTodoEndpoint))
		handler.ServeHTTP(w, req)
	})
	http.HandleFunc("/deleteTodo", func(w http.ResponseWriter, req *http.Request) {
		var handler http.Handler = corsHandler.Handler(http.HandlerFunc(deleteTodoEndpoint))
		handler.ServeHTTP(w, req)
	})

	http.HandleFunc("/editTodo", func(w http.ResponseWriter, req *http.Request) {
		var handler http.Handler = corsHandler.Handler(http.HandlerFunc(editTodoEndpoint))
		handler.ServeHTTP(w, req)
	})
	http.ListenAndServe("localhost:"+port, nil)
}