package notification

import "github.com/byliuyang/kgs/app/entity"

type Notifier interface {
	NotifyRequester(message entity.Message, requester entity.Requester) error
}
