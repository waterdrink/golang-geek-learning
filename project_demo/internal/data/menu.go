package data

import (
	"context"
	"learning/project_demo/internal/biz"
	"log"
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

func (m *MenuRepo) GetMenus(ctx context.Context) (menus []*biz.Menu, err error) {
	for _, m := range m.menu {
		log.Printf("repo: get menu %v\n", m.Id)
		menus = append(menus, m)
	}
	return menus, nil
}

func (m *MenuRepo) SaveMenu(ctx context.Context, menu *biz.Menu) (err error) {
	log.Printf("repo: save menu %v\n", menu.Id)
	m.menu[menu.Id] = menu
	return nil
}
