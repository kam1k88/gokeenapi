package cmd

import (
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newConfigureAwgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "configure-awg",
		Aliases: []string{"configureawg"},
		Short:   "Configure Wireguard connection to add ASC parameters to it in Keenetic router",
	}

	cmd.Flags().String("interface-id", "", "Keenetic interface ID to configure")
	cmd.Flags().String("conf-file", "", "Path to a conf TOML file with AWG configuration")

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag(config.ViperKeeneticInterfaceId, cmd.Flags().Lookup("interface-id"))
		_ = viper.BindPFlag(config.ViperKeeneticInterfaceConfFile, cmd.Flags().Lookup("conf-file"))
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return gokeenrestapi.AwgConf.ConfigureOrUpdateInterface(viper.GetString(config.ViperKeeneticInterfaceConfFile), viper.GetString(config.ViperKeeneticInterfaceId))
	}
	return cmd
}
