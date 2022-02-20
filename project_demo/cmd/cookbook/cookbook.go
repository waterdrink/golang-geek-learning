package main

import (
	"context"
	"fmt"
	"learning/project_demo/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	echo "github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

type cookbookApp struct {
	Echo        *echo.Echo
	MenuService *service.MenuService
}

var myCookbookApp *cookbookApp

const RouterPrefix = "/api/v1/cookbook"

func newCookbookApp(menuService *service.MenuService) *cookbookApp {
	app := &cookbookApp{
		Echo:        echo.New(),
		MenuService: menuService,
	}
	app.Echo.Group(RouterPrefix)
	app.Echo.GET("/menu", app.MenuService.GetMenuV1)
	app.Echo.POST("/menu", app.MenuService.SaveMenuV1)
	return app
}

func main() {
	gctx, cancel := context.WithCancel(context.Background())
	g, errCtx := errgroup.WithContext(gctx)

	myCookbookApp = initCookbookApp()

	// start http server
	g.Go(func() error {
		err := myCookbookApp.Echo.Start(":9009")
		if nil != err && http.ErrServerClosed != err {
			return err
		}
		return nil
	})

	// stop http server when group context done
	g.Go(func() error {
		<-errCtx.Done()
		fmt.Printf("shutdown http server...\n")
		// shutdown timeout
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		return myCookbookApp.Echo.Shutdown(ctx)
	})

	// handle exit signal
	g.Go(func() error {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-errCtx.Done():
			return fmt.Errorf("unexpected shutdown because: %v", errCtx.Err())
		case <-exit:
			fmt.Println("receive exit signal")
			cancel()
			return nil
		}

	})

	if err := g.Wait(); nil != err {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("shutdown success !")
}
