package generator

import (
	"os"
	"strings"
	"text/template"
)

func getDir(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join(parts[:len(parts)-1], "/")
}

func generateFile(tmplPath, outPath string, data any) error {
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	tmpl, err := template.New("tmpl").
		Funcs(template.FuncMap{
			"lower": strings.ToLower,
		}).Parse(string(content))
	if err != nil {
		return err
	}

	_ = os.MkdirAll(getDir(outPath), os.ModePerm)
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return tmpl.Execute(out, data)
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
