package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	gctx, cancel := context.WithCancel(context.Background())
	g, errCtx := errgroup.WithContext(gctx)

	srv := http.Server{
		Addr:    ":9009",
		Handler: http.DefaultServeMux,
	}

	// start http server
	g.Go(func() error {
		http.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request) {
			fmt.Println("hello !")
		})

		err := srv.ListenAndServe()
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
		return srv.Shutdown(ctx)
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
