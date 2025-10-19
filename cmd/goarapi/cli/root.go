package cli

import (
	"strings"

	"github.com/fatih/color"
	"github.com/kam1k88/gokeenapi/internal/gokeenlog"
	"github.com/kam1k88/gokeenapi/internal/gokeenversion"
	"github.com/kam1k88/gokeenapi/pkg/config"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{}
	rootCmd.SilenceUsage = true
	var configFile string
	rootCmd.Use = "gokeenapi"
	rootCmd.SilenceErrors = true
	rootCmd.Short = "Automate your Keenetic router management with simple commands"
	rootCmd.Long = `gokeenapi - Automate your Keenetic router management with ease

A powerful command-line utility for managing Keenetic routers via REST API.
Supports route management, DNS configuration, WireGuard setup, and more.

Key features:
• Add/delete static routes from .bat files and URLs
• Manage DNS records for local domain resolution  
• Configure WireGuard (AWG) VPN connections
• Clean up known hosts with pattern matching
• Execute custom router commands directly
• Works with both local IP and KeenDNS addresses

Examples:
  # Show all available interfaces
  gokeenapi show-interfaces --config config.yaml

  # Add routes from configuration
  gokeenapi add-routes --config config.yaml

  # Set up WireGuard connection
  gokeenapi add-awg --config config.yaml --conf-file vpn.conf

For detailed command help, use: gokeenapi [command] --help`

	rootCmd.PersistentFlags().Bool("debug", false,
		`Enable debug mode with verbose logging.
Shows detailed API requests, responses, and internal operations.
Useful for troubleshooting connection or configuration issues.`)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "",
		`Path to YAML configuration file (required).
Contains router connection details and operation settings.
Can also be set via GOKEENAPI_CONFIG environment variable.`)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// completion and help commands should run without any checks and init
		commandsToSkip := []string{CmdCompletion, CmdHelp}
		for _, commandToSkip := range commandsToSkip {
			if strings.Contains(cmd.CommandPath(), commandToSkip) {
				return nil
			}
		}
		err := config.LoadConfig(configFile)
		if err != nil {
			return err
		}
		if err = checkRequiredFields(); err != nil {
			return err
		}
		gokeenlog.Infof("🚀  %v: %v, %v: %v", color.BlueString("Version"), color.CyanString(gokeenversion.Version()), color.BlueString("Build date"), color.CyanString(gokeenversion.BuildDate()))
		gokeenlog.Info("🏗️  Configuration loaded:")
		gokeenlog.InfoSubStepf("%v: %v", color.BlueString("Keenetic URL"), config.Cfg.Keenetic.URL)
		gokeenlog.InfoSubStepf("%v: %v", color.BlueString("Config"), color.CyanString(configFile))
		return prepareAPI(cmd.Context())
	}

	rootCmd.AddCommand(
		newAddRoutesCmd(),
		newDeleteRoutesCmd(),
		newShowInterfacesCmd(),
		newUpdateAwgCmd(),
		newAddAwgCmd(),
		newAddDnsRecordsCmd(),
		newDeleteDnsRecordsCmd(),
		newDeleteKnownHostsCmd(),
		newExecCmd(),
		newServeCmd(),
	)
	return rootCmd
}
