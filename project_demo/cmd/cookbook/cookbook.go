package main

import (
	"context"
	"fmt"
	"learning/project_demo/internal/pkg/mdns"
	"learning/project_demo/internal/rpcservice"
	"learning/project_demo/internal/service"
	pb "learning/project_demo/proto"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	echo "github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type cookbookApp struct {
	ID             int // 微服务唯一ID
	Echo           *echo.Echo
	mdnsServer     *mdns.MdnsService // 使用mdns服务发现
	MenuService    *service.MenuService
	MenuRpcService *rpcservice.MenuRpcService
}

var myCookbookApp *cookbookApp

const RouterPrefix = "/api/v1/cookbook"

func newCookbookApp(menuService *service.MenuService, rpcService *rpcservice.MenuRpcService) *cookbookApp {
	rand.Seed(time.Now().Unix())
	app := &cookbookApp{
		ID:             rand.Int(),
		Echo:           echo.New(),
		MenuService:    menuService,
		MenuRpcService: rpcService,
	}
	app.Echo.Group(RouterPrefix)
	app.Echo.GET("/menu", app.MenuService.GetMenuV1)
	app.Echo.POST("/menu", app.MenuService.SaveMenuV1)
	log.Printf("cookbook app id: %d\n", app.ID)
	return app
}

func main() {
	gctx, cancel := context.WithCancel(context.Background())
	g, errCtx := errgroup.WithContext(gctx)

	myCookbookApp = initCookbookApp()

	// start rpc server
	g.Go(func() error {
		lis, err := net.Listen("tcp", ":9008")
		if err != nil {
			return fmt.Errorf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterMenuServer(grpcServer, myCookbookApp.MenuRpcService)
		return grpcServer.Serve(lis)
	})

	// start http server
	g.Go(func() error {
		err := myCookbookApp.Echo.Start(":9009")
		if nil != err && http.ErrServerClosed != err {
			return err
		}
		return nil
	})

	// listen mdns
	// 开启服务发现
	g.Go(func() error {
		var err error
		myCookbookApp.mdnsServer, err = mdns.NewMdnsService(strconv.Itoa(myCookbookApp.ID), "cookbook", "eth0")
		if nil != err {
			return err
		}
		myCookbookApp.MenuService.RegisterServiceDiscover(myCookbookApp.mdnsServer)

		return myCookbookApp.mdnsServer.StartDiscover()
	})

	// stop http server when group context done
	g.Go(func() error {
		<-errCtx.Done()
		fmt.Printf("shutdown http server...\n")
		// shutdown timeout
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		err := myCookbookApp.Echo.Shutdown(ctx)
		fmt.Printf("shutdown http server done\n")
		return err
	})

	// stop mdns server when group context done
	g.Go(func() error {
		<-errCtx.Done()
		log.Printf("shutdown mdns server...\n")
		if nil != myCookbookApp.mdnsServer {
			return myCookbookApp.mdnsServer.Shutdown()
		}
		log.Printf("shutdown mdns server done\n")
		return nil
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
