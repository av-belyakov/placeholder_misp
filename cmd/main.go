package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	go func() {
		log.Printf("system call:%+v", <-ctx.Done())

		stop()
	}()

	server(ctx)
}
