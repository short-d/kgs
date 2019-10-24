package message

import (
	"time"

	"github.com/byliuyang/kgs/app/adapter/template"
	"github.com/byliuyang/kgs/app/entity"
)

func NewKeyGenSucceedMessage(
	tmpl template.Template,
	timeElapsed time.Duration,
) (entity.Message, error) {
	body, err := tmpl.Render("key-gen-succeed.gohtml",
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
