package cmd

import (
	"github.com/joho/godotenv"
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/internal/keenlog"
	"github.com/spf13/cobra"
	"os"
	"strings"
)
import "github.com/spf13/viper"

func NewRootCmd() *cobra.Command {
	viper.Reset()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.MustBindEnv(config.ViperKeeneticApi, "GOKEENAPI_API")
	viper.MustBindEnv(config.ViperKeeneticLogin, "GOKEENAPI_LOGIN")
	viper.MustBindEnv(config.ViperKeeneticPassword, "GOKEENAPI_PASSWORD")
	viper.MustBindEnv(config.ViperDebug, "GOKEENAPI_DEBUG")
	viper.MustBindEnv(config.ViperKeeneticInterfaceId, "GOKEENAPI_INTERFACE")
	viper.MustBindEnv(config.ViperKeeneticConfig, "GOKEENAPI_CONFIG")
	rootCmd := &cobra.Command{}
	rootCmd.SilenceUsage = true
	var configFile string
	rootCmd.Use = "gokeenapi"
	rootCmd.Short = "A convenient utility to add/delete routes in Keenetic routers via REST API"
	rootCmd.PersistentFlags().String("api", "", "Keenetic API url, should contain /rci at the end. Example: https://api.my-super-keenetic.keenetic.pro/rci")
	rootCmd.PersistentFlags().String("login", "", "Keenetic API login")
	rootCmd.PersistentFlags().String("password", "", "Keenetic API password")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode and logging")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to YAML config file (optional)")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		_ = viper.BindPFlag(config.ViperDebug, rootCmd.PersistentFlags().Lookup("debug"))
		_ = viper.BindPFlag(config.ViperKeeneticApi, rootCmd.PersistentFlags().Lookup("api"))
		_ = viper.BindPFlag(config.ViperKeeneticLogin, rootCmd.PersistentFlags().Lookup("login"))
		_ = viper.BindPFlag(config.ViperKeeneticPassword, rootCmd.PersistentFlags().Lookup("password"))
		_ = viper.BindPFlag(config.ViperKeeneticConfig, rootCmd.PersistentFlags().Lookup("config"))
		if configFile != "" {
			viper.SetConfigFile(configFile)
			viper.SetConfigType("yaml")
			err := viper.ReadInConfig()
			if err != nil {
				return err
			}
		}
		_, statErr := os.Stat(".gokeenapienv")
		if statErr == nil {
			err := godotenv.Load(".gokeenapienv")
			if err != nil {
				return err
			}
		}
		err := checkRequiredFields()
		if err != nil {
			return err
		}
		keenlog.Infof("API to use: %v", viper.GetString(config.ViperKeeneticApi))
		return nil
	}

	rootCmd.AddCommand(
		newAddRoutesCmd(),
		newDeleteRoutesCmd(),
		newShowInterfacesCmd(),
	)
	return rootCmd
}
