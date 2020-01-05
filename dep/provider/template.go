package provider

import (
	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/app/modern/mdtemplate"
)

type TemplateRootDir string

func NewHTML(rootDir TemplateRootDir) fw.Template {
	return mdtemplate.NewHTML(string(rootDir))
}
