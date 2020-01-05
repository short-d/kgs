package template

import "time"

type KeyGenSucceedData struct {
	TimeElapsed time.Duration
}

var KeyGenSucceedTemplate = "key-gen-succeed.gohtml"

var KeyGenSucceedIncludeTemplates = []string{
	"key-gen-succeed.gohtml",
}
