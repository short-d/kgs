package provider

import "github.com/byliuyang/kgs/app/adapter/template"

type TemplatePattern string

func NewTemplate(templatePattern TemplatePattern) (template.Template, error) {
	return template.NewTemplate(string(templatePattern))
}
