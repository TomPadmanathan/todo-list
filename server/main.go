package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
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

func writeToFile(filename string, data []byte) error {
	// open/create json file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// write to file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func readFromFile(filename string) error {
	// Open file or create it if it doesn't exist
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read from file
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// If file is empty, store an empty array
	if len(data) == 0 {
		emptyArray := []byte("[]")
		_, err = file.Write(emptyArray)
		if err != nil {
			return err
		}
		_, err = file.Seek(0, 0) // Move the file cursor back to the beginning
		if err != nil {
			return err
		}
		todos = []todo{}
		return nil
	}

	// Unmarshal file
	err = json.Unmarshal(data, &todos)
	if err != nil {
		return err
	}

	return nil
}


func todosEndpoint(w http.ResponseWriter, req *http.Request) {
	if(req.Method != "GET") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// read/update todos from json
	var err error = readFromFile("todos.json")
	if err != nil {
		http.Error(w, "error reading todos from json", http.StatusInternalServerError)
		return
	}

	responseTodos(w)
}

func addTodoEndpoint(w http.ResponseWriter, req *http.Request) {
	if(req.Method != "POST") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get newtodo data
	var body, err = io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "todo not found", http.StatusBadRequest)
		return
	}
	// construct new todo
	var newTodo todo

	newTodo.Todo = string(body)
	newTodo.Id = uuid.New().String()
	newTodo.Timestamp = time.Now().Unix()

	// add new todo 
	todos = append(todos, newTodo)

	updateTodosInJSONFile(w)
	
	responseTodos(w)
}

func updateTodosInJSONFile(w http.ResponseWriter)  {
	// convert todos to json format
	jsonData, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		http.Error(w, "error writing todos to json", http.StatusInternalServerError)
		return
	}

	// write the JSON data to file
	err = writeToFile("todos.json", jsonData)
	if err != nil {
		http.Error(w, "error writing todos to json", http.StatusInternalServerError)
		return
	}

}

func deleteTodoEndpoint(w http.ResponseWriter, req *http.Request) {
	if(req.Method != "POST") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get todoId
	var todoId, err = io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "todoId not found", http.StatusBadRequest)
		return
	}

	// get the index to delete
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

	// delete todo
	todos = append(todos[:deleteItemindex], todos[deleteItemindex+1:]...)

	updateTodosInJSONFile(w)

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

	updateTodosInJSONFile(w)

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