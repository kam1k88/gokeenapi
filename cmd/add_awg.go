package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAddAwgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-awg",
		Aliases: []string{"addawg"},
		Short:   "Add AWG connection from config file in Keenetic router",
	}

	cmd.Flags().String("conf-file", "", "Path to a conf TOML file with WG configuration")
	var configure bool
	var up bool
	var name string
	cmd.Flags().BoolVar(&configure, "configure", false, "Add ASC parameters to the connection after creating from config file")
	cmd.Flags().BoolVar(&up, "up", false, "Bring interface up after creating")
	cmd.Flags().StringVar(&name, "name", "", "Name for new WG interface, optional")
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag(config.ViperKeeneticInterfaceConfFile, cmd.Flags().Lookup("conf-file"))
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		confPath := viper.GetString(config.ViperKeeneticInterfaceConfFile)
		if confPath == "" {
			return fmt.Errorf("conf-file flag is required")
		}
		confPath, err := filepath.Abs(confPath)
		if err != nil {
			return err
		}
		err = gokeenrestapi.Checks.CheckAWGInterfaceExistsFromConfFile(confPath)
		if err != nil {
			return err
		}
		gokeenlog.InfoSubStepf("Conf-file: %v", color.CyanString("%v", confPath))
		createdInterface, err := gokeenrestapi.AwgConf.AddInterface(confPath, name)
		if err != nil {
			return err
		}
		gokeenlog.Info("Created Wireguard interface!")
		gokeenlog.InfoSubStepf("Id: %v", color.CyanString(createdInterface.Created))
		if configure {
			err = gokeenrestapi.AwgConf.ConfigureOrUpdateInterface(confPath, createdInterface.Created)
		}
		if err != nil {
			return err
		}
		err = gokeenrestapi.Interface.SetGlobalIpInInterface(createdInterface.Created, true)
		if err != nil {
			return err
		}
		if up {
			err = gokeenrestapi.Interface.UpInterface(createdInterface.Created)
			if err != nil {
				return err
			}
			err = gokeenrestapi.Interface.WaitUntilInterfaceIsUp(createdInterface.Created)
		}
		return err
	}
	return cmd
}
