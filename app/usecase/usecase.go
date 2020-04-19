package usecase

import (
	"fmt"
	"time"

	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/app/entity"
	"github.com/short-d/kgs/app/usecase/keys"
	"github.com/short-d/kgs/app/usecase/notification"
)

type UseCase struct {
	logger          fw.Logger
	producer        keys.Producer
	consumer        keys.Consumer
	eventDispatcher fw.Emitter
}

func (u UseCase) PopulateKey(keyLength uint, requesterEmail string) {
	startAt := time.Now()

	if err := u.producer.Produce(keyLength); err != nil {
		u.logger.Error(err)
		return
	}

	err := u.eventDispatcher.Dispatch(notification.OnKeyPopulatedEvent{
		TimeElapsed: time.Since(startAt),
		Requester: entity.Requester{
			Name:  "",
			Email: requesterEmail,
		},
	})

	if err != nil {
		u.logger.Error(err)
		return
	}

	u.logger.Info("Finish populating keys")
}

func (u UseCase) AllocateKeys(maxKeyCount uint) ([]string, error) {
	allocatedKeys, err := u.consumer.ConsumeInBatch(maxKeyCount)
	if err != nil {
		return nil, err
	}
	u.logger.Info(fmt.Sprintf("Allocated %d keys to client", len(allocatedKeys)))
	return allocatedKeys, nil
}

func NewUseCase(
	logger fw.Logger,
	producer keys.Producer,
	consumer keys.Consumer,
	eventDispatcher fw.Emitter,
) UseCase {
	return UseCase{
		logger:          logger,
		producer:        producer,
		consumer:        consumer,
		eventDispatcher: eventDispatcher,
	}
}
