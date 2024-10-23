package main

import (
	"context"
	"fmt"
	"github.com/noksa/gokeenapi/cmd"
	"log"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	ctx := context.Background()
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to execute command: %v", err))
	}
}
