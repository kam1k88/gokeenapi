package cmd

import (
	"strings"

	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/internal/gokeenversion"
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{}
	rootCmd.SilenceUsage = true
	var configFile string
	rootCmd.Use = "gokeenapi"
	rootCmd.SilenceErrors = true
	rootCmd.Short = "A utility to run commands (such as add/del routes) in Keenetic routers via REST API"
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode and logging")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to YAML config file (required)")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// completion and help commands should run without any checks and init
		commandsToSkip := []string{"completion", "help"}
		for _, commandToSkip := range commandsToSkip {
			if strings.Contains(cmd.CommandPath(), commandToSkip) {
				return nil
			}
		}
		err := config.LoadConfig(configFile)
		if err != nil {
			return err
		}
		err = checkRequiredFields()
		if err != nil {
			return err
		}
		gokeenlog.Infof("Version: %v, BuildDate: %v", gokeenversion.Version(), gokeenversion.BuildDate())
		gokeenlog.Info("Configuration loaded:")
		gokeenlog.InfoSubStepf("Keenetic url: %v", config.Cfg.Keenetic.URL)
		gokeenlog.InfoSubStepf("Config: %v", configFile)
		return gokeenrestapi.Common.Auth()
	}

	rootCmd.AddCommand(
		newAddRoutesCmd(),
		newDeleteRoutesCmd(),
		newShowInterfacesCmd(),
		newUpdateAwgCmd(),
		newAddAwgCmd(),
		newAddDnsRecordsCmd(),
		newDeleteDnsRecordsCmd(),
	)
	return rootCmd
}
