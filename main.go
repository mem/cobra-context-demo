package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mem/cobra-context-demo/cmd"
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	defer func() {
		fmt.Println("Ending...")
		signal.Stop(signals)
		cancel()
	}()

	go func() {
		select {
		case <-signals:
			fmt.Println("Got signal, propagating...")
			cancel()

		case <-ctx.Done():
		}
	}()

	if err := cmd.Execute(ctx); err != nil {
		fmt.Println(err)
		exitCode = 1
		return
	}
}
