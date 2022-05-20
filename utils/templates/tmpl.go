package templates

import (
	"html/template"
	"io"
)

type Template struct {
	t *template.Template
}

func (t Template) Execute(writer io.Writer, data string) error {
	return t.t.Execute(writer, template.HTML(data))
}

func (t Template) Any(writer io.Writer, data any) error {
	return t.t.Execute(writer, data)
}

func (t Template) Internal(writer io.Writer) error {
	return t.t.Execute(writer, "internal error, try again")
}

func New() *Templates {
	return new(Templates)
}
