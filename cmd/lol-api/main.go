package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/real-web-world/lol-api/bootstrap"
	"github.com/real-web-world/lol-api/global"
)

func main() {
	begin := time.Now()
	engine, err := bootstrap.InitApp()
	defer global.Cleanup()
	if err != nil {
		panic(fmt.Sprintf("初始化应用失败:%v\n", err))
	}
	httpAddr := fmt.Sprintf("%s:%d", global.Conf.HTTPHost,
		global.Conf.HTTPPort)
	srv := &http.Server{
		Addr:    httpAddr,
		Handler: engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("初始化app成功,env:%s,耗时%v\n", global.Conf.Mode, time.Since(begin))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Println("server shutdown err:", err)
		return
	}
	log.Println("server shutdown success")
}
