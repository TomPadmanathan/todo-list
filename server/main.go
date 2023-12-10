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
		http.Error(w, "reviewId not found", http.StatusBadRequest)
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

	var reviewId, err = io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "reviewId not found", http.StatusBadRequest)
		return
	}
	var deleteItemindex int
	for i := 0; i < len(todos); i++ {
		if(todos[i].Id == string(reviewId)) {
			deleteItemindex = i
			break
		}
	}
	if deleteItemindex < 0 || deleteItemindex >= len(todos) {
		http.Error(w, "invalid delete index", http.StatusBadRequest)
    }
	todos = append(todos[:deleteItemindex], todos[deleteItemindex+1:]...)

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
	http.ListenAndServe("localhost:"+port, nil)
}