package service

import (
	"fmt"
	"learning/project_demo/internal/biz"
	pkgError "learning/project_demo/internal/pkg/error"
	"net/http"

	"github.com/labstack/echo/v4"
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
	Menu MenuRespV1 `json:"menu" description:"menu"`
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
	menu, err := m.mc.GetMenu(c.Request().Context(), req.Id)
	if nil != err {
		return NewGenericErrResp(c, err, pkgError.Internal)
	}
	return c.JSON(http.StatusOK, GetMenuRespV1{
		Menu: MenuRespV1{
			Id:   menu.Id,
			Name: menu.Name,
		},
	})
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
