package provider

import (
	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/app/usecase/notification"
)

func NewEventEmitter(
	eventDispatcher fw.Dispatcher,
	emailNotifierListener notification.EmailNotifierEventListener,
) (fw.Emitter, error) {
	err := eventDispatcher.BindListeners([]fw.Listener{
		emailNotifierListener,
	})

	if err != nil {
		return nil, err
	}

	return eventDispatcher, nil
}
