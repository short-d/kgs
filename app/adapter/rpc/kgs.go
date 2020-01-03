package rpc

import (
	"context"

	"github.com/byliuyang/kgs/app/adapter/rpc/proto"
	"github.com/byliuyang/kgs/app/usecase"
	"github.com/golang/protobuf/ptypes/empty"
)

var _ proto.KeyGenServer = (*KeyGenServer)(nil)

type KeyGenServer struct {
	useCase usecase.UseCase
}

func (k KeyGenServer) AllocateKeys(
	ctx context.Context,
	req *proto.AllocateKeysRequest,
) (*proto.AllocateKeysResponse, error) {
	allocatedKeys, err := k.useCase.AllocateKeys(uint(req.MaxKeyCount))
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
		k.useCase.PopulateKey(uint(req.KeyLength), req.RequesterEmail)
	}()
	return &empty.Empty{}, nil
}

func NewKeyGenServer(
	useCase usecase.UseCase,
) KeyGenServer {
	return KeyGenServer{
		useCase: useCase,
	}
}
