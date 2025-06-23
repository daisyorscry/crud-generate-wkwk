package generator

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"lower": strings.ToLower,
	"sqlType": func(goType string) string {
		switch goType {
		case "string":
			return "VARCHAR(255)"
		case "text":
			return "TEXT"
		case "bool":
			return "BOOLEAN"
		case "int":
			return "INTEGER"
		case "float64":
			return "DOUBLE PRECISION"
		default:
			return "TEXT"
		}
	},
}
