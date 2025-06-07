
package generator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GenerateModelRequest struct {
	Model  string  `json:"model"`
	Fields []Field `json:"fields"`
}

func HandleGenerateModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range req.Fields {
		req.Fields[i].Name = strings.Title(req.Fields[i].Name)
		req.Fields[i].Type = mapType(req.Fields[i].Type)
		req.Fields[i].Tag = fmt.Sprintf("`gorm:\"column:%s\" json:\"%s\"`", strings.ToLower(req.Fields[i].Name), strings.ToLower(req.Fields[i].Name))
	}

	data := map[string]any{
		"Model":  strings.Title(req.Model),
		"Fields": req.Fields,
	}

	outPath := fmt.Sprintf("models/%s.go", strings.ToLower(req.Model))
	err := generateFile("templates/model.go.tmpl", outPath, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Model %s generated at %s\n", req.Model, outPath)
}
