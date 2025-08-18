package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/cmd"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	ctx := context.Background()
	cmdCtx, _ := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	var err error
	errChan := make(chan error)
	go func() {
		errChan <- rootCmd.ExecuteContext(cmdCtx)
	}()
	select {
	case err = <-errChan:
	case <-cmdCtx.Done():
		err = cmdCtx.Err()
	}
	cmd.RestoreCursor()
	if err != nil {
		gokeenlog.Info(color.RedString("Error occured!"))
		os.Exit(100)
	}
}
