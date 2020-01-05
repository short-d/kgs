package provider

import "github.com/short-d/app/modern/mdemail"

type SendGridAPIKey string

func NewSendGrid(apiKey SendGridAPIKey) mdemail.SendGrid {
	return mdemail.NewSendGrid(string(apiKey))
}
