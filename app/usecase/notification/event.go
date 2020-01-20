package notification

import (
	"time"

	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/app/entity"
)

const onKeyPopulatedEventName = "onKeyPopulatedEvent"

var _ fw.Event = (*OnKeyPopulatedEvent)(nil)

type OnKeyPopulatedEvent struct {
	TimeElapsed time.Duration
	Requester   entity.Requester
}

func (e OnKeyPopulatedEvent) GetName() string {
	return onKeyPopulatedEventName
}
