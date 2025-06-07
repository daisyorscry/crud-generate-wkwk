package pkgErr

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type FieldErrors map[string]string

type ErrorDetail struct {
	Detail error `json:"-"`
}

func Wrap(err error) *ErrorDetail {
	if err == nil {
		return nil
	}
	return &ErrorDetail{Detail: err}
}

func (e *ErrorDetail) MarshalJSON() ([]byte, error) {
	if e == nil || e.Detail == nil {
		return []byte(`"null"`), nil
	}
	escaped := strings.ReplaceAll(e.Detail.Error(), `"`, `\"`)
	return []byte(`"` + escaped + `"`), nil
}

func NewErrorDetail(err error) *ErrorDetail {
	if err == nil {
		return nil
	}
	return &ErrorDetail{Detail: err}
}

var defaultTemplates = map[string]string{
	// Common
	"required": "{field} is required",
	"email":    "{field} must be a valid email address",
	"min":      "{field} must be at least {param} characters",
	"max":      "{field} must be at most {param} characters",
	"len":      "{field} must be exactly {param} characters",
	"eq":       "{field} must be equal to {param}",
	"ne":       "{field} must not be equal to {param}",
	"gt":       "{field} must be greater than {param}",
	"gte":      "{field} must be greater than or equal to {param}",
	"lt":       "{field} must be less than {param}",
	"lte":      "{field} must be less than or equal to {param}",
	"eqfield":  "{field} must be equal to {param}",
	"nefield":  "{field} must not be equal to {param}",
	"oneof":    "{field} must be one of: {param}",

	// String
	"alpha":           "{field} must contain only letters",
	"alphanum":        "{field} must contain only letters and numbers",
	"alphanumunicode": "{field} must contain only letters and numbers (unicode)",
	"ascii":           "{field} must contain only ASCII characters",
	"boolean":         "{field} must be a boolean value",
	"contains":        "{field} must contain '{param}'",
	"containsany":     "{field} must contain any of '{param}'",
	"excludes":        "{field} must not contain '{param}'",
	"lowercase":       "{field} must be all lowercase",
	"uppercase":       "{field} must be all uppercase",
	"startswith":      "{field} must start with '{param}'",
	"endswith":        "{field} must end with '{param}'",

	// Format
	"uuid":        "{field} must be a valid UUID",
	"uuid4":       "{field} must be a valid UUIDv4",
	"datetime":    "{field} must be a valid datetime",
	"json":        "{field} must be valid JSON",
	"credit_card": "{field} must be a valid credit card number",
	"base64":      "{field} must be a valid base64 string",
	"e164":        "{field} must be a valid E.164 phone number",
	"url":         "{field} must be a valid URL",
	"http_url":    "{field} must be a valid HTTP(S) URL",
	"ip":          "{field} must be a valid IP address",
	"ipv4":        "{field} must be a valid IPv4 address",
	"ipv6":        "{field} must be a valid IPv6 address",
	"mac":         "{field} must be a valid MAC address",

	// File / System
	"file":  "{field} must be a valid file",
	"dir":   "{field} must be a valid directory",
	"image": "{field} must be an image file",

	// Special
	"isdefault":   "{field} must be default value",
	"excluded_if": "{field} must not be set when condition is met",
	"unique":      "{field} must be unique",
}

func ParseValidationErrors(err error) FieldErrors {
	errors := FieldErrors{}
	if err == nil {
		return errors
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := getJSONFieldName(fe)
			tag := fe.Tag()
			param := fe.Param()
			template := defaultTemplates[tag]

			// fallback message
			if template == "" {
				if param != "" {
					template = "{field} failed on rule '" + tag + "' with param '" + param + "'"
				} else {
					template = "{field} failed on rule '" + tag + "'"
				}
			}

			msg := strings.ReplaceAll(template, "{field}", humanizeFieldName(field))
			msg = strings.ReplaceAll(msg, "{param}", param)

			errors[field] = msg
		}
	}
	return errors
}

func getJSONFieldName(fe validator.FieldError) string {
	if field, ok := reflect.TypeOf(fe).FieldByName(fe.StructField()); ok {
		if tag := field.Tag.Get("json"); tag != "" {
			return strings.Split(tag, ",")[0]
		}
	}
	return strings.ToLower(fe.Field())
}

var titleCaser = cases.Title(language.English)

func humanizeFieldName(name string) string {
	parts := strings.Split(name, "_")
	for i, p := range parts {
		parts[i] = titleCaser.String(p)
	}
	return strings.Join(parts, " ")
}
