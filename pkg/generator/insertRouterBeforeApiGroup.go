package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func insertAfterApiGroup(outputPath, templatePath string, data any) error {
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
		newLines = append(newLines, line)
		if !inserted && strings.Contains(line, `api := app.Group("/api")`) {
			// Insert langsung setelah baris api := app.Group("/api")
			newLines = append(newLines, insert)
			inserted = true
		}
	}
	if !inserted {
		return fmt.Errorf("could not find `api := app.Group(\"/api\")` in %s", outputPath)
	}
	return os.WriteFile(outputPath, []byte(strings.Join(newLines, "\n")), 0644)
}
