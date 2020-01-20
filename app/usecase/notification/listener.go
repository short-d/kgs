package notification

import (
	"errors"
	"fmt"

	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/app/adapter/message"
)

var _ fw.Listener = (*EmailNotifierEventListener)(nil)

type EmailNotifierEventListener struct {
	logger              fw.Logger
	template            fw.Template
	serviceName         string
	serviceEmailAddress string
	emailSender         fw.EmailSender
}

func (l EmailNotifierEventListener) GetSubscribedEvent() string {
	return onKeyPopulatedEventName
}

func (l EmailNotifierEventListener) Handle(event fw.Event) {
	e, ok := event.(OnKeyPopulatedEvent)

	if !ok {
		l.logger.Error(errors.New("expected OnKeyPopulatedEvent"))
		return
	}

	msg, err := message.NewKeyGenSucceedMessage(l.template, e.TimeElapsed)

	if err != nil {
		l.logger.Error(fmt.Errorf("failed to generated succeed message: %v", err))
		return
	}

	email := fw.Email{
		FromName:    l.serviceName,
		FromAddress: l.serviceEmailAddress,
		ToName:      e.Requester.Name,
		ToAddress:   e.Requester.Email,
		Subject:     msg.Title,
		ContentHTML: msg.BodyHTML,
	}

	if err := l.emailSender.SendEmail(email); err != nil {
		l.logger.Error(fmt.Errorf("failed to send email: %v", err))
		return
	}
}

func NewEmailNotifierEventListener(
	logger fw.Logger,
	template fw.Template,
	serviceName string,
	serviceEmailAddress string,
	emailSender fw.EmailSender,
) EmailNotifierEventListener {
	return EmailNotifierEventListener{
		logger:              logger,
		template:            template,
		serviceName:         serviceName,
		serviceEmailAddress: serviceEmailAddress,
		emailSender:         emailSender,
	}
}
