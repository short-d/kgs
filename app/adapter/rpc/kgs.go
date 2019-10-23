package rpc

import (
	"context"

	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/app/usecase/keys"
	"github.com/golang/protobuf/ptypes/empty"
)

var _ KeyGenServer = (*KeyGenController)(nil)

type KeyGenController struct {
	producer keys.Producer
	consumer keys.Consumer
	logger   fw.Logger
}

func (k KeyGenController) AllocateKeys(
	ctx context.Context,
	req *AllocateKeysRequest,
) (*AllocateKeysResponse, error) {
	allocatedKeys, err := k.consumer.ConsumeInBatch(uint(req.MaxKeyCount))
	if err != nil {
		return &AllocateKeysResponse{}, err
	}
	return &AllocateKeysResponse{Keys: allocatedKeys}, nil
}

func (k KeyGenController) PopulateKeys(
	ctx context.Context,
	req *PopulateKeysRequest,
) (*empty.Empty, error) {
	go func() {
		err := k.producer.Produce(uint(req.KeyLength))
		if err != nil {
			k.logger.Error(err)
		}
	}()
	return &empty.Empty{}, nil
}

func NewKeyGenController(
	producer keys.Producer,
	consumer keys.Consumer,
	logger fw.Logger,
) KeyGenController {
	return KeyGenController{
		producer: producer,
		consumer: consumer,
		logger:   logger,
	}
}
