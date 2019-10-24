package notification

import (
	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/app/entity"
)

var _ Notifier = (*EmailNotifier)(nil)

type EmailNotifier struct {
	serviceName         string
	serviceEmailAddress string
	emailSender         fw.EmailSender
}

func (e EmailNotifier) NotifyRequester(
	message entity.Message,
	requester entity.Requester,
) error {
	email := fw.Email{
		FromName:    e.serviceName,
		FromAddress: e.serviceEmailAddress,
		ToName:      requester.Name,
		ToAddress:   requester.Email,
		Subject:     message.Title,
		ContentHTML: message.BodyHTML,
	}

	return e.emailSender.SendEmail(email)
}

func NewEmailNotifier(
	serviceName string,
	serviceEmailAddress string,
	emailSender fw.EmailSender,
) EmailNotifier {
	return EmailNotifier{
		serviceName:         serviceName,
		serviceEmailAddress: serviceEmailAddress,
		emailSender:         emailSender,
	}
}
