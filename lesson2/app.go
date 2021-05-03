package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type App struct{
	ctx context.Context
	cancel func()
	servers []ServerInterface
}

func NewApp(servers ...ServerInterface) *App{
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx: ctx,
		cancel: cancel,
		servers: servers,
	}
}

func (a *App) Run() error{
	g, ctx := errgroup.WithContext(a.ctx)
	for _,srv := range a.servers{
		srv := srv
		g.Go(func()error{
			<-ctx.Done()
			fmt.Println("receive server stop")
			return srv.Stop()
		})
		g.Go(func()error{
			fmt.Println("server start")
			return srv.Start()
		})
	}
	c:=make(chan os.Signal,1)
	signal.Notify(c,syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	g.Go(func()error{
		for{
			select {
			case <-c:
				fmt.Println("app stop")
				a.Stop()
			case <-ctx.Done():
				fmt.Println("go exit")
				return ctx.Err()
			}
		}
	})
	if err := g.Wait();err!= nil && !errors.Is(err,context.Canceled){
		return err
	}
	return nil
}

func(a *App)Stop() {
	if a.cancel!= nil{
		a.cancel()
	}
}