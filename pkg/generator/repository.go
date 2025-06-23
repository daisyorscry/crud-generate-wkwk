package generator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GenerateRepoRequest struct {
	Model string `json:"model"`
}

func HandleGenerateRepository(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRepoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := strings.ToTitle(req.Model)
	modelVar := strings.ToLower(req.Model)

	data := map[string]any{
		"Model":    model,
		"modelVar": modelVar,
	}

	outPath := fmt.Sprintf("repositories/%s_repository.go", modelVar)
	err := generateFile("templates/repository.go.tmpl", outPath, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Repository %s generated at %s\n", model, outPath)
}
