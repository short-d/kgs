package message

import (
	"time"

	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/app/adapter/template"
	"github.com/short-d/kgs/app/entity"
)

func NewKeyGenSucceedMessage(
	tmpl fw.Template,
	timeElapsed time.Duration,
) (entity.Message, error) {
	body, err := tmpl.Render(
		template.KeyGenSucceedTemplate,
		template.KeyGenSucceedIncludeTemplates,
		template.KeyGenSucceedData{
			TimeElapsed: timeElapsed,
		})
	if err != nil {
		return entity.Message{}, err
	}
	return entity.Message{
		Title:    "Key Gen Status Update",
		BodyHTML: body,
	}, nil
}
