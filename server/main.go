package main

import (
	"encoding/json"
	"net/http"

	"github.com/rs/cors"
)

func todos(w http.ResponseWriter, req *http.Request) {
	type todo struct {
		Id string `json:"id"`
		Todo string `json:"todo"`
		Timestamp string `json:"timestamp"`
	}
	var todos []todo = []todo{
		{Id: "1", Todo: "Wash dishes", Timestamp: "1702020402"},
		{Id: "2", Todo: "Get Landary", Timestamp: "1701020402"},
	}
	if(req.Method != "GET") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jsonData, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	const port string = "3001"

	corsHandler := cors.Default()
	handler := corsHandler.Handler(http.HandlerFunc(todos))
	http.HandleFunc("/todos", func(w http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(w, req)
	})
	http.ListenAndServe("localhost:"+port, nil)
}