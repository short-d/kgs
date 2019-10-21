package rpc

import (
	"github.com/byliuyang/app/fw"
	"google.golang.org/grpc"
)

var _ fw.GRpcAPI = (*KgsAPI)(nil)

type KgsAPI struct {
	keyGenServer KeyGenServer
}

func (k KgsAPI) RegisterServers(server *grpc.Server) {
	RegisterKeyGenServer(server, k.keyGenServer)
}

func NewKgsAPI(keyGenServer KeyGenServer) KgsAPI {
	return KgsAPI{keyGenServer: keyGenServer}
}
