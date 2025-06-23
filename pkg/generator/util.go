package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// func getDir(path string) string {
// 	parts := strings.Split(path, "/")
// 	return strings.Join(parts[:len(parts)-1], "/")
// }

func generateFile(templatePath string, outputPath string, data any) error {
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// pastikan folder tujuan ada
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func mapType(input string) string {
	switch strings.ToLower(input) {
	case "string":
		return "string"
	case "int":
		return "int"
	case "uint":
		return "uint"
	case "bool":
		return "bool"
	case "float":
		return "float64"
	case "text":
		return "string"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}
