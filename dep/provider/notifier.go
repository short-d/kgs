package provider

import (
	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/app/usecase/notification"
)

type ServiceEmailAddress string

func NewEmailNotifier(
	serviceName string,
	serviceEmailAddress ServiceEmailAddress,
	emailSender fw.EmailSender,
) notification.EmailNotifier {
	return notification.NewEmailNotifier(
		serviceName,
		string(serviceEmailAddress),
		emailSender,
	)
}
