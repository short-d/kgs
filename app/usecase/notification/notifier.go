package notification

import "github.com/short-d/kgs/app/entity"

type Notifier interface {
	NotifyRequester(message entity.Message, requester entity.Requester) error
}
