package rpc

import (
	"context"

	"github.com/byliuyang/kgs/app/usecase/keys/producer"
	"github.com/golang/protobuf/ptypes/empty"
)

var _ KeyGenServer = (*KeyGenController)(nil)

type KeyGenController struct {
	producer producer.Producer
}

func (k KeyGenController) PopulateKeys(
	ctx context.Context,
	req *PopulateKeysRequest,
) (*empty.Empty, error) {
	go func() {
		k.producer.Produce(uint(req.KeySize))
	}()
	return &empty.Empty{}, nil
}

func NewKeyGenController(producer producer.Producer) KeyGenController {
	return KeyGenController{
		producer: producer,
	}
}
