package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/kam1k88/gokeenapi/cmd/goarapi/cli"
	"github.com/kam1k88/gokeenapi/internal/gokeenlog"
)

func main() {
	rootCmd := cli.NewRootCmd()
	ctx := context.Background()
	cmdCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errChan := make(chan error)
	go func() {
		errChan <- rootCmd.ExecuteContext(cmdCtx)
	}()

	var err error
	select {
	case err = <-errChan:
	case <-cmdCtx.Done():
		err = cmdCtx.Err()
	}

	cli.RestoreCursor()
	if err != nil {
		gokeenlog.Info(color.HiRedString("Error occured!"))
		gokeenlog.Info(color.RedString(err.Error()))
		os.Exit(100)
	}
}
