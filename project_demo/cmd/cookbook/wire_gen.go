// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"learning/project_demo/internal/biz"
	"learning/project_demo/internal/data"
	"learning/project_demo/internal/service"
)

// Injectors from wire.go:

// initCookbookApp init my cookbook application.
func initCookbookApp() *cookbookApp {
	menuRepo := data.NewMenuRepo()
	menuUsecase := biz.NewMenuUsecase(menuRepo)
	menuService := service.NewMenuService(menuUsecase)
	mainCookbookApp := newCookbookApp(menuService)
	return mainCookbookApp
}
