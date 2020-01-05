package provider

import (
	"github.com/short-d/app/fw"
	"github.com/short-d/app/modern/mdtemplate"
)

type TemplateRootDir string

func NewHTML(rootDir TemplateRootDir) fw.Template {
	return mdtemplate.NewHTML(string(rootDir))
}
