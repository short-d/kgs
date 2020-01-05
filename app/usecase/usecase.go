package usecase

import (
	"fmt"
	"time"

	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/app/adapter/message"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/keys"
	"github.com/byliuyang/kgs/app/usecase/notification"
)

type UseCase struct {
	logger   fw.Logger
	template fw.Template
	producer keys.Producer
	consumer keys.Consumer
	notifier notification.Notifier
}

func (u UseCase) PopulateKey(keyLength uint, requesterEmail string) {
	startAt := time.Now()
	err := u.producer.Produce(keyLength)
	if err != nil {
		u.logger.Error(err)
		return
	}

	timeElapsed := time.Since(startAt)
	msg, err := message.NewKeyGenSucceedMessage(u.template, timeElapsed)
	if err != nil {
		u.logger.Error(err)
		return
	}

	requester := entity.Requester{
		Name:  "",
		Email: requesterEmail,
	}
	err = u.notifier.NotifyRequester(msg, requester)
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
	template fw.Template,
	producer keys.Producer,
	consumer keys.Consumer,
	notifier notification.Notifier,
) UseCase {
	return UseCase{
		logger:   logger,
		template: template,
		producer: producer,
		consumer: consumer,
		notifier: notifier,
	}
}
