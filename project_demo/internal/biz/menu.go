package biz

import (
	"context"
)

type Menu struct {
	Id   int
	Name string
}

type MenuRepo interface {
	GetMenus(ctx context.Context) (menu []*Menu, err error)
	SaveMenu(ctx context.Context, menu *Menu) (err error)
}

type MenuUsecase struct {
	MenuRepo MenuRepo
}

func NewMenuUsecase(menuRepo MenuRepo) *MenuUsecase {
	return &MenuUsecase{
		MenuRepo: menuRepo,
	}
}

func (m *MenuUsecase) GetMenus(ctx context.Context) (menu []*Menu, err error) {
	return m.MenuRepo.GetMenus(ctx)
}

func (m *MenuUsecase) SaveMenu(ctx context.Context, menu *Menu) (err error) {
	return m.MenuRepo.SaveMenu(ctx, menu)
}
