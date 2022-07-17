package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/**
退出程序之后，console 打印的结果如下：
2022/07/17 23:48:32 errgroup exit
2022/07/17 23:48:32 server shutting down
errgroup exiting: os signal: interrupt
*/
func main() {

	eg, ctx := errgroup.WithContext(context.Background())

	mut := http.NewServeMux()
	mut.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, geek"))
	})

	//模拟单个服务错误
	quitServer := make(chan struct{})
	mut.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		quitServer <- struct{}{}
	})

	server := http.Server{
		Handler: mut,
		Addr:    ":8080",
	}

	//启动http server
	eg.Go(func() error {
		return server.ListenAndServe()
	})

	//监听服务是否异常，如果异常shut down
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("errgroup exit")
		case <-quitServer:
			log.Println("server quit")
		}
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		log.Println("server shutting down")
		return server.Shutdown(timeoutCtx)
	})

	//监听系统的信号，退出所有协成
	eg.Go(func() error {

		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return fmt.Errorf("os signal: %v ", sig)
		}
	})

	fmt.Printf("errgroup exiting: %+v\n", eg.Wait())

}
