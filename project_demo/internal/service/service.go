package service

import (
	"learning/project_demo/internal/biz"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewMenuService)

type MenuService struct {
	mc *biz.MenuUsecase
}

// @title cookbook menu app
// @version 1.0
// @description This is a cookbook menu app.
func NewMenuService(mc *biz.MenuUsecase) *MenuService {
	return &MenuService{
		mc: mc,
	}
}
