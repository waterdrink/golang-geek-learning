package data

import (
	"context"
	"learning/project_demo/internal/biz"
)

var _ biz.MenuRepo = (*MenuRepo)(nil)

type MenuRepo struct {
	menu map[int]*biz.Menu
}

func NewMenuRepo() biz.MenuRepo {
	return &MenuRepo{
		menu: make(map[int]*biz.Menu),
	}
}

func (m *MenuRepo) GetMenu(ctx context.Context, id int) (menu *biz.Menu, err error) {
	return m.menu[id], nil
}

func (m *MenuRepo) SaveMenu(ctx context.Context, menu *biz.Menu) (err error) {
	m.menu[menu.Id] = menu
	return nil
}
