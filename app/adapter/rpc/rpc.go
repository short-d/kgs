package rpc

import (
	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/app/adapter/rpc/proto"
	"google.golang.org/grpc"
)

var _ fw.GRpcAPI = (*KgsAPI)(nil)

type KgsAPI struct {
	keyGenServer proto.KeyGenServer
}

func (k KgsAPI) RegisterServers(server *grpc.Server) {
	proto.RegisterKeyGenServer(server, k.keyGenServer)
}

func NewKgsAPI(keyGenServer proto.KeyGenServer) KgsAPI {
	return KgsAPI{keyGenServer: keyGenServer}
}
