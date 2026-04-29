package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/XenomorphingTV/waymark/exporter"
	"github.com/XenomorphingTV/waymark/parser"
)

func main() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	http.HandleFunc("/api/story", handleStory)

	fmt.Println("visualizer running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("server error: ", err)
		os.Exit(1)
	}
}

func handleStory(w http.ResponseWriter, r *http.Request) {
	// get filename from query par
	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "missing file parameter", http.StatusBadRequest)
		return
	}

	src, err := os.ReadFile(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read file: %v", err), http.StatusNotFound)
		return
	}

	tokens, err := parser.Tokenize(string(src))
	if err != nil {
		http.Error(w, fmt.Sprintf("tokenize error: %v", err), http.StatusUnprocessableEntity)
		return
	}

	story, err := parser.Parse(tokens)
	if err != nil {
		http.Error(w, fmt.Sprintf("parse error: %v", err), http.StatusUnprocessableEntity)
		return
	}

	data, err := exporter.Export(story)
	if err != nil {
		http.Error(w, fmt.Sprintf("export error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
