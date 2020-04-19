package provider

import (
	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/app/usecase/notification"
)

type ServiceEmailAddress string

func NewEmailNotifierEventListener(
	logger fw.Logger,
	template fw.Template,
	serviceName string,
	serviceEmailAddress ServiceEmailAddress,
	emailSender fw.EmailSender,
) notification.EmailNotifierEventListener {
	return notification.NewEmailNotifierEventListener(
		logger,
		template,
		serviceName,
		string(serviceEmailAddress),
		emailSender,
	)
}
