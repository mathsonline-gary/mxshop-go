const (
{{ range .Errors }}
Code{{ .CamelValue }} ErrorCode = {{ .Number }}
{{- end }}
)

{{ range .Errors }}

func IsError{{.CamelValue}}(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
    return e.Code == int32(Code{{ .CamelValue }})
}

func Error{{ .CamelValue }}(format string, args ...interface{}) *errors.Error {
	 return errors.New(int32(Code{{ .CamelValue }}), fmt.Sprintf(format, args...))
}

{{- end }}