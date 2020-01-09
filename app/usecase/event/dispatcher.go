package event

import (
	"errors"
	"sync"

	"github.com/asaskevich/EventBus"
)

var _ Dispatcher = (*EventDispatcher)(nil)
var _ Subscriber = (*EventDispatcher)(nil)

// ErrDispatcherIsClosed tells that there is no way to perform manipulations with event dispatcher
var ErrDispatcherIsClosed = errors.New("failed to perform the operation, the dispatcher is closed")

// EventDispatcher publishes an event to all its subscribers
type EventDispatcher struct {
	bus      EventBus.Bus
	lock     sync.RWMutex
	isClosed bool
}

func (d *EventDispatcher) Dispatch(event Event) error {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if d.isClosed {
		return ErrDispatcherIsClosed
	}

	d.bus.Publish(event.GetName(), event)

	return nil
}

func (d *EventDispatcher) Subscribe(eventName string, listener Listener) error {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if d.isClosed {
		return ErrDispatcherIsClosed
	}

	return d.bus.SubscribeAsync(eventName, listener.Handle, false)
}

func (d *EventDispatcher) Unsubscribe(eventName string, listener Listener) error {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if d.isClosed {
		return ErrDispatcherIsClosed
	}

	return d.bus.Unsubscribe(eventName, listener.Handle)
}

func (d *EventDispatcher) Close() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.isClosed {
		return ErrDispatcherIsClosed
	}

	d.isClosed = true
	d.bus.WaitAsync()

	return nil
}

// NewEventDispatcher creates a new instance of EventDispatcher type
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		bus:      EventBus.New(),
		lock:     sync.RWMutex{},
		isClosed: false,
	}
}
