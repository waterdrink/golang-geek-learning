//+build wireinject

package main

import (
	"learning/project_demo/internal/biz"
	"learning/project_demo/internal/data"
	"learning/project_demo/internal/service"

	"github.com/google/wire"
)

// initCookbookApp init my cookbook application.
func initCookbookApp() *cookbookApp {
	wire.Build(newCookbookApp, service.ProviderSet, biz.ProviderSet, data.ProviderSet)
	return &cookbookApp{}
}
