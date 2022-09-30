package main

import (
	"context"
	"flag"
	"github.com/PavelDonchenko/bookstoreCRUD/cmd/server"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	go func() {
		server.Run()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	<-ctx.Done()

	log.Println("shutting down")
	os.Exit(0)
}
