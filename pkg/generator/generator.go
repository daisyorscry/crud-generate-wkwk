package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

	if strings.TrimSpace(req.Feature) == "" {
		http.Error(w, "feature is required", http.StatusBadRequest)
		return
	}

	feature := strings.ToLower(req.Feature)
	model := toTitle(req.Model)
	modelVar := strings.ToLower(req.Model)

	// format fields
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

	base := fmt.Sprintf("domain/%s", feature)

	// Tambah inisialisasi sebelum return d
	if err := insertBeforeReturn("domain/domain.go", "templates/domain_init.go.tmpl", data); err != nil {
		http.Error(w, "Failed to update domain.go => "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Tambah field ke struct Domain
	if err := appendStructField("domain/domain.go", fmt.Sprintf("%sHandler *%sRest.%sHandler", model, modelVar, model)); err != nil {
		http.Error(w, "Failed to append to Domain struct => "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Tambah import otomatis jika belum ada
	if err := appendImportsIfMissing("domain/domain.go", []string{
		fmt.Sprintf("\"daisy/domain/%s/repository\"", feature),
		fmt.Sprintf("\"daisy/domain/%s/service\"", feature),
		fmt.Sprintf("\"daisy/domain/%s/handler/rest\"", feature), // ⬅ tetap ada, tapi akan diberi alias
		"\"github.com/go-playground/validator/v10\"",
	}); err != nil {
		http.Error(w, "Failed to append imports => "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := insertAfterApiGroup("routes/route.go", "templates/route.go.tmpl", data); err != nil {
		http.Error(w, "Failed to update routes/route.go => "+err.Error(), http.StatusInternalServerError)
		return
	}

	files := []struct {
		Template string
		Output   string
	}{
		{"templates/create_table.sql.tmpl", fmt.Sprintf("%s/sql/create_%s_table.sql", base, modelVar)},

		{"templates/dao_request.go.tmpl", fmt.Sprintf("%s/dao/request.go", base)},
		{"templates/dao_response.go.tmpl", fmt.Sprintf("%s/dao/response.go", base)},
		{"templates/dao.go.tmpl", fmt.Sprintf("%s/dao/dao.go", base)},
		{"templates/mapper.go.tmpl", fmt.Sprintf("%s/dao/mapper.go", base)},

		{"templates/repository.go.tmpl", fmt.Sprintf("%s/repository/base.go", base)},
		{"templates/service.go.tmpl", fmt.Sprintf("%s/service/%s_service.go", base, modelVar)},
		{"templates/rest.go.tmpl", fmt.Sprintf("%s/handler/rest/%s_handler.go", base, modelVar)},
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

func insertBeforeReturn(outputPath, templatePath string, data any) error {
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	insert := buf.String()

	contentBytes, err := os.ReadFile(outputPath)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	if strings.Contains(content, insert) {
		return nil // sudah ada
	}

	lines := strings.Split(content, "\n")
	var newLines []string
	inserted := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inserted && strings.HasPrefix(trimmed, "return") && strings.Contains(trimmed, "d") {
			newLines = append(newLines, strings.TrimRight(insert, "\n"))
			inserted = true
		}
		newLines = append(newLines, line)
	}

	if !inserted {
		return fmt.Errorf("could not find `return d` in %s", outputPath)
	}

	return os.WriteFile(outputPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func appendFile(outputPath string, templatePath string, data any) error {
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// baca isi file dulu
	content, _ := os.ReadFile(outputPath)
	if bytes.Contains(content, buf.Bytes()) {
		return nil // sudah ada
	}

	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	return err
}

func appendStructField(filePath, fieldLine string) error {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	// Cek apakah field sudah ada
	if strings.Contains(content, fieldLine) {
		return nil
	}

	lines := strings.Split(content, "\n")
	var newLines []string
	inserted := false
	inStruct := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		newLines = append(newLines, line)

		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "type Domain struct {") {
			inStruct = true
			continue
		}

		if inStruct && trimmed == "}" {
			// Sisipkan sebelum tanda kurung tutup
			newLines = append(newLines[:len(newLines)-1], "\t"+fieldLine, "}")
			inserted = true
			inStruct = false
			continue
		}
	}

	if !inserted {
		return fmt.Errorf("could not find Domain struct closing brace in %s", filePath)
	}

	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

func appendImportsIfMissing(outputPath string, newImports []string) error {
	contentBytes, err := os.ReadFile(outputPath)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	startIdx := strings.Index(content, "import (")
	if startIdx == -1 {
		return fmt.Errorf("import block not found in %s", outputPath)
	}
	endIdx := strings.Index(content[startIdx:], ")")
	if endIdx == -1 {
		return fmt.Errorf("import block not closed in %s", outputPath)
	}
	endIdx += startIdx

	importBlock := content[startIdx : endIdx+1]

	for _, imp := range newImports {
		// pakai alias kalau path mengandung '/handler/', '/service/', '/repository/'
		var alias string
		if strings.Contains(imp, "/handler/") {
			alias = getAliasFromPath(imp, "Rest")
		} else if strings.Contains(imp, "/repository") {
			alias = getAliasFromPath(imp, "Repo")
		} else if strings.Contains(imp, "/service") {
			alias = getAliasFromPath(imp, "Service")
		}

		fullImport := fmt.Sprintf("%s \"%s\"", alias, strings.Trim(imp, "\""))
		if !strings.Contains(importBlock, fullImport) {
			importBlock = strings.Replace(importBlock, ")", fmt.Sprintf("\t%s\n)", fullImport), 1)
		}
	}

	newContent := content[:startIdx] + importBlock + content[endIdx+1:]
	return os.WriteFile(outputPath, []byte(newContent), 0644)
}

// helper
func getAliasFromPath(importPath, suffix string) string {
	parts := strings.Split(importPath, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "domain" && i+1 < len(parts) {
			return strings.ToLower(parts[i+1]) + suffix
		}
	}
	return suffix
}
