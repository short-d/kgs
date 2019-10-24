package template

import (
	"bytes"
	"html/template"
)

type Template struct {
	template *template.Template
}

func (t Template) Render(templateFileName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := t.template.ExecuteTemplate(&buf, templateFileName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func NewTemplate(templatePattern string) (Template, error) {
	tmpl, err := template.ParseGlob(templatePattern)
	if err != nil {
		return Template{}, err
	}
	return Template{
		template: tmpl,
	}, nil
}
