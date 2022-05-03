package service

import (
	"context"
	"fmt"
	"learning/project_demo/internal/biz"
	pkgError "learning/project_demo/internal/pkg/error"
	"learning/project_demo/proto"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GenericErrResp struct {
	ErrorMsg  string             `json:"error_msg"`
	ErrorCode pkgError.ErrorCode `json:"error_code"`
}

func NewGenericErrResp(c echo.Context, err error, code pkgError.ErrorCode) error {
	return c.JSON(http.StatusOK, GenericErrResp{ErrorMsg: fmt.Sprintf("%v", err), ErrorCode: code})
}

func NewGenericOKResp(c echo.Context) error {
	return c.JSON(http.StatusOK, GenericErrResp{ErrorMsg: "", ErrorCode: pkgError.OK})
}

type GetMenuReqV1 struct {
	Id int `json:"id" description:"menu id"`
}

type GetMenuRespV1 struct {
	Menus []MenuRespV1 `json:"menu" description:"menu"`
}

type MenuRespV1 struct {
	Id   int    `json:"id" description:"menu id"`
	Name string `json:"name" description:"menu name"`
}

// GetMenuV1
// @Summary get menu
// @Description get menu
// @Tags menu
// @Id GetMenuV1
// @Accept json
// @Produce json
// @Param id query GetMenuReqV1 true "menu id"
// @Success 200 {object} GetMenuRespV1
// @Failure default {object} GenericErrResp
// @router /menu [get]
func (m *MenuService) GetMenuV1(c echo.Context) (err error) {
	req := new(GetMenuReqV1)
	if err = c.Bind(req); err != nil {
		return NewGenericErrResp(c, err, pkgError.InvalidArgument)
	}
	menus, err := m.mc.GetMenus(c.Request().Context())
	if nil != err {
		return NewGenericErrResp(c, err, pkgError.Internal)
	}
	otherMenus, err := m.getMenusFromOtherNode()
	if nil != err {
		return NewGenericErrResp(c, err, pkgError.Internal)
	}

	respMenus := []MenuRespV1{}
	for _, menu := range append(menus, otherMenus...) {
		respMenus = append(respMenus, MenuRespV1{
			Id:   menu.Id,
			Name: menu.Name,
		})
	}

	return c.JSON(http.StatusOK, GetMenuRespV1{Menus: respMenus})
}

func (m *MenuService) getMenusFromOtherNode() (menus []*biz.Menu, err error) {
	if m.mdnsServer == nil {
		return nil, fmt.Errorf("mdns server is nil")
	}

	for _, addr := range m.mdnsServer.GetDiscoveredServiceAddr() {
		log.Printf("start rpc call to %s\n", addr)
		// Set up a connection to the server.
		conn, err := grpc.Dial(fmt.Sprintf("%v:9008", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("did not connect: %v\n", err)
			continue
		}
		defer conn.Close()
		c := proto.NewMenuClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		output, err := c.GetMenus(ctx, &proto.GetMenusInput{})
		if err != nil {
			log.Printf("could not get menus from %v: %v\n", addr, err)
			continue
		}
		for _, m := range output.GetMenus() {
			log.Printf("got menus from %v: %v", addr, m.MenuId)
			menus = append(menus, &biz.Menu{
				Id:   int(m.MenuId),
				Name: m.MenuName,
			})
		}
	}
	return menus, nil
}

type SaveMenuReqV1 struct {
	Menu MenuRespV1 `json:"menu" description:"menu"`
}

// SaveMenuV1
// @Summary save menu
// @Description save menu
// @Tags menu
// @Id SaveMenuV1
// @Accept json
// @Produce json
// @Param menu query SaveMenuReqV1 true "menu"
// @Success 200 {object} GenericErrResp
// @Failure default {object} GenericErrResp
// @router /menu [post]
func (m *MenuService) SaveMenuV1(c echo.Context) (err error) {
	req := new(SaveMenuReqV1)
	if err = c.Bind(req); err != nil {
		return NewGenericErrResp(c, err, pkgError.InvalidArgument)
	}
	if err := m.mc.SaveMenu(c.Request().Context(), &biz.Menu{
		Id:   req.Menu.Id,
		Name: req.Menu.Name,
	}); nil != err {
		return NewGenericErrResp(c, err, pkgError.Internal)
	}
	return NewGenericOKResp(c)
}
