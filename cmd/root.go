package cmd

import (
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/spf13/cobra"
	"strings"
)
import "github.com/spf13/viper"

func NewRootCmd() *cobra.Command {
	viper.Reset()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.MustBindEnv(config.ViperKeeneticApi, "KEENETIC_API")
	viper.MustBindEnv(config.ViperKeeneticLogin, "KEENETIC_LOGIN")
	viper.MustBindEnv(config.ViperKeeneticPassword, "KEENETIC_PASSWORD")
	viper.MustBindEnv(config.ViperDebug, "KEENETIC_DEBUG")
	viper.MustBindEnv(config.ViperKeeneticInterfaceId, "KEENETIC_INTERFACE")
	rootCmd := &cobra.Command{}
	rootCmd.SilenceUsage = true
	rootCmd.Use = "A convenient utility to work with Keenetic router via REST API"
	rootCmd.PersistentFlags().String("api", "", "Keenetic API url, should contain /rci at the end. Example: https://api.my-super-keenetic.keenetic.pro/rci")
	rootCmd.PersistentFlags().String("login", "", "Keenetic API login")
	rootCmd.PersistentFlags().String("password", "", "Keenetic API password")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode and logging")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		_ = viper.BindPFlag(config.ViperDebug, rootCmd.PersistentFlags().Lookup("debug"))
		_ = viper.BindPFlag(config.ViperKeeneticApi, rootCmd.PersistentFlags().Lookup("api"))
		_ = viper.BindPFlag(config.ViperKeeneticLogin, rootCmd.PersistentFlags().Lookup("login"))
		_ = viper.BindPFlag(config.ViperKeeneticPassword, rootCmd.PersistentFlags().Lookup("password"))
		return checkRequiredFields()
	}

	rootCmd.AddCommand(
		newAddRoutesCmd(),
		newDeleteRoutesCmd(),
		newShowInterfacesCmd(),
	)
	return rootCmd
}
