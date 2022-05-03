package rpcservice

import (
	"context"
	pb "learning/project_demo/proto"
)

func (m *MenuRpcService) GetMenus(c context.Context, _ *pb.GetMenusInput) (*pb.GetMenusOutput, error) {
	menus, err := m.mc.GetMenus(c)
	if nil != err {
		return nil, err
	}
	respMenus := []*pb.MenuInfo{}
	for _, menu := range menus {
		respMenus = append(respMenus, &pb.MenuInfo{
			MenuId:   int64(menu.Id),
			MenuName: menu.Name,
		})
	}

	return &pb.GetMenusOutput{Menus: respMenus}, nil
}
