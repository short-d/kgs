package provider

import (
	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/app/usecase/notification"
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
