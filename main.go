package main

import (
	"daisy/pkg/generator"
	"net/http"
)

func main() {
	http.HandleFunc("/generate-model", generator.HandleGenerateModel)
	http.HandleFunc("/generate-repository", generator.HandleGenerateRepository)
	http.HandleFunc("/generate-all", func(w http.ResponseWriter, r *http.Request) {
		generator.HandleGenerateAll(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
