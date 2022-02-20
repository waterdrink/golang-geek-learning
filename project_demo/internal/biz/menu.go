package biz

import "context"

type Menu struct {
	Id   int
	Name string
}

type MenuRepo interface {
	GetMenu(ctx context.Context, id int) (menu *Menu, err error)
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

func (m *MenuUsecase) GetMenu(ctx context.Context, id int) (menu *Menu, err error) {
	return m.MenuRepo.GetMenu(ctx, id)
}

func (m *MenuUsecase) SaveMenu(ctx context.Context, menu *Menu) (err error) {
	return m.MenuRepo.SaveMenu(ctx, menu)
}
