package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kam1k88/gokeenapi/internal/gokeenlog"
	"github.com/kam1k88/gokeenapi/pkg/server"
	"github.com/spf13/cobra"
)

func newServeCmd() *cobra.Command {
	var addr string
	cmd := &cobra.Command{
		Use:     CmdServe,
		Aliases: AliasesServe,
		Short:   "Expose REST API for router management",
		RunE: func(cmd *cobra.Command, args []string) error {
			api := API()
			if api == nil {
				return fmt.Errorf("api is not initialized")
			}
			srv := server.New(api)
			gokeenlog.Infof("🌐  Serving REST API on %s", color.CyanString(addr))
			return srv.ListenAndServe(addr)
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":8080", "HTTP listen address")
	return cmd
}
