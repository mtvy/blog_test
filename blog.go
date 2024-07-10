package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type List struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ID struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	lists = []List{
		{ID: 1, Title: "First Blog", Content: "This is the first blog post."},
		{ID: 2, Title: "Second Blog", Content: "This is the second blog post."},
	}
	ids = []ID{
		{ID: 1, Name: "Author One"},
		{ID: 2, Name: "Author Two"},
	}
	mutex = &sync.Mutex{}
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/blog/list", blogListHandler)
	http.HandleFunc("/blog/id", blogIDHandler)

	log.Println("Server is listening on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func blogListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST,DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case http.MethodGet:
		handleGetList(w)
	case http.MethodPost:
		handlePostList(w, r)
	case http.MethodDelete:
		handleDelList(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetList(w http.ResponseWriter) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	jsonData, err := json.Marshal(lists)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func handlePostList(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var newList List
	if err := json.Unmarshal(body, &newList); err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	newList.ID = len(lists) + 1
	lists = append(lists, newList)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func handleDelList(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	idStr := strings.TrimPrefix(r.URL.Path, "/blog/list/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, list := range lists {
		if list.ID == id {
			lists = append(lists[:i], lists[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "List not found", http.StatusNotFound)

}

func blogIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case http.MethodGet:
		handleGetID(w)
	case http.MethodPost:
		handlePostID(w, r)
	case http.MethodDelete:
		handleDelID(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetID(w http.ResponseWriter) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	jsonData, err := json.Marshal(ids)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func handlePostID(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST,DELETE,OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var newID ID
	if err := json.Unmarshal(body, &newID); err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	newID.ID = len(ids) + 1
	ids = append(ids, newID)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func handleDelID(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	idStr := strings.TrimPrefix(r.URL.Path, "/blog/id/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, ident := range ids {
		if ident.ID == id {
			ids = append(ids[:i], ids[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "ID not found", http.StatusNotFound)
}
