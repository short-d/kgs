package rpc

import (
	"context"
	"time"

	"github.com/byliuyang/kgs/app/adapter/rpc/proto"

	"github.com/byliuyang/kgs/app/adapter/message"
	"github.com/byliuyang/kgs/app/adapter/template"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/notification"

	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/app/usecase/keys"
	"github.com/golang/protobuf/ptypes/empty"
)

var _ proto.KeyGenServer = (*KeyGenServer)(nil)

type KeyGenServer struct {
	producer keys.Producer
	consumer keys.Consumer
	notifier notification.Notifier
	template template.Template
	logger   fw.Logger
}

func (k KeyGenServer) AllocateKeys(
	ctx context.Context,
	req *proto.AllocateKeysRequest,
) (*proto.AllocateKeysResponse, error) {
	allocatedKeys, err := k.consumer.ConsumeInBatch(uint(req.MaxKeyCount))
	if err != nil {
		return &proto.AllocateKeysResponse{}, err
	}
	return &proto.AllocateKeysResponse{Keys: allocatedKeys}, nil
}

func (k KeyGenServer) PopulateKeys(
	ctx context.Context,
	req *proto.PopulateKeysRequest,
) (*empty.Empty, error) {
	go func() {
		startAt := time.Now()
		k.logger.Info("Start populating keys")
		err := k.producer.Produce(uint(req.KeyLength))
		if err != nil {
			k.logger.Error(err)
			return
		}

		timeElapsed := time.Now().Sub(startAt)
		msg, err := message.NewKeyGenSucceedMessage(k.template, timeElapsed)
		if err != nil {
			k.logger.Error(err)
			return
		}

		requester := entity.Requester{
			Name:  "",
			Email: req.RequesterEmail,
		}

		err = k.notifier.NotifyRequester(msg, requester)
		if err != nil {
			k.logger.Error(err)
			return
		}
		k.logger.Info("Finish populating keys")
	}()
	return &empty.Empty{}, nil
}

func NewKeyGenServer(
	producer keys.Producer,
	consumer keys.Consumer,
	notifier notification.Notifier,
	template template.Template,
	logger fw.Logger,
) KeyGenServer {
	return KeyGenServer{
		producer: producer,
		consumer: consumer,
		notifier: notifier,
		template: template,
		logger:   logger,
	}
}
