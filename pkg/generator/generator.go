package generator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GenerateAllRequest struct {
	Feature string  `json:"feature"`
	Model   string  `json:"model"`
	Fields  []Field `json:"fields"`
}

var titleCaser = cases.Title(language.English)

func toTitle(s string) string {
	return titleCaser.String(s)
}

func HandleGenerateAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateAllRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	feature := strings.ToLower(req.Feature)
	model := toTitle(req.Model)
	modelVar := strings.ToLower(req.Model)

	// format field types & tags
	for i := range req.Fields {
		req.Fields[i].Name = toTitle(req.Fields[i].Name)
		req.Fields[i].Type = mapType(req.Fields[i].Type)
		req.Fields[i].Tag = fmt.Sprintf("`json:\"%s\" validate:\"required\"`", strings.ToLower(req.Fields[i].Name))
	}

	data := map[string]any{
		"Feature":  feature,
		"Model":    model,
		"modelVar": modelVar,
		"Fields":   req.Fields,
	}

	if strings.TrimSpace(feature) == "" {
		http.Error(w, "feature is required", http.StatusBadRequest)
		return
	}

	base := fmt.Sprintf("domain/%s", feature)

	files := []struct {
		Template string
		Output   string
	}{
		{"templates/dao_request.go.tmpl", fmt.Sprintf("%s/dao/request.go", base)},
		{"templates/dao_response.go.tmpl", fmt.Sprintf("%s/dao/response.go", base)},
		{"templates/dao.go.tmpl", fmt.Sprintf("%s/dao/dao.go", base)},
		{"templates/repository.go.tmpl", fmt.Sprintf("%s/repository/%s_repository.go", base, modelVar)},
		{"templates/service.go.tmpl", fmt.Sprintf("%s/service/%s_service.go", base, modelVar)},
		{"templates/rest.go.tmpl", fmt.Sprintf("%s/handler/rest/%s_handler.go", base, modelVar)},
		{"templates/mapper.go.tmpl", fmt.Sprintf("%s/dao/mapper.go", base)},
	}

	for _, f := range files {
		if err := generateFile(f.Template, f.Output, data); err != nil {
			http.Error(w, "Failed to generate: "+f.Output+" => "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "All components for feature '%s' generated successfully", feature)
}
