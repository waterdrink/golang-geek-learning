package rpcservice

import (
	"learning/project_demo/internal/biz"
	pb "learning/project_demo/proto"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewMenuRpcService)

type MenuRpcService struct {
	mc *biz.MenuUsecase
	pb.UnimplementedMenuServer
}

func NewMenuRpcService(mc *biz.MenuUsecase) *MenuRpcService {
	return &MenuRpcService{
		mc: mc,
	}
}
