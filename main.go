package main

import (
	"context"
	"github.com/noksa/gokeenapi/cmd"
	"os"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	ctx := context.Background()
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(100)
	}
}
