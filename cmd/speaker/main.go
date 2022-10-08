package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/exitstop/speaker_alpine/internal/console"
	"github.com/exitstop/speaker_alpine/internal/translateshell"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		cancel()
		os.Exit(0)
	}()

	trShell := translateshell.New(ctx)
	go trShell.Run()

	//go func() {
	console.Add(cancel, trShell)
	console.Low()
	//}()

}
