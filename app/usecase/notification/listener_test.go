package notification

import (
	"testing"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/short-d/app/fw"
	"github.com/short-d/app/mdtest"
	"github.com/short-d/app/modern/mdevent"
	"github.com/short-d/kgs/app/entity"
)

func TestEmailEventHandle(t *testing.T) {
	logger := mdtest.NewLoggerFake(mdtest.FakeLoggerArgs{})
	sender := fakeEmailSender{}

	listener := NewEmailNotifierEventListener(
		&logger,
		fakeTemplator{},
		"serviceName",
		"serviceEmailAddress",
		&sender,
	)

	dispatcher := mdevent.NewEventDispatcher(EventBus.New())
	mdtest.Equal(t, nil, dispatcher.Subscribe(listener))

	err := dispatcher.Dispatch(OnKeyPopulatedEvent{
		TimeElapsed: time.Second,
		Requester: entity.Requester{
			Name:  "",
			Email: "recipientEmailAddress",
		},
	})

	mdtest.Equal(t, nil, err)
	mdtest.Equal(t, nil, dispatcher.Close())

	expected := fw.Email{
		FromName:    "serviceName",
		FromAddress: "serviceEmailAddress",
		ToName:      "",
		ToAddress:   "recipientEmailAddress",
		Subject:     "Key Gen Status Update",
		ContentHTML: "contentHTML",
	}

	mdtest.Equal(t, expected, sender.email)
}

type fakeTemplator struct{}

func (t fakeTemplator) Render(renderTemplate string, includeTemplates []string, data interface{}) (string, error) {
	return "contentHTML", nil
}

type fakeEmailSender struct {
	email fw.Email
}

func (s *fakeEmailSender) SendEmail(email fw.Email) error {
	s.email = email

	return nil
}
