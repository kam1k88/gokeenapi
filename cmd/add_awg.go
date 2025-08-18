package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/internal/gokeenspinner"
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
	cmd.Flags().String("interface-name", "", "Name for new WG interface, optional")
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
		err = checkAWGInterfaceExistsFromConfFile(confPath)
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
			err = configureInterface(confPath, createdInterface.Created)
		}
		if err != nil {
			return err
		}
		if up {
			err = gokeenrestapi.Interface.UpInterface(createdInterface.Created)
			if err != nil {
				return err
			}
			err = gokeenspinner.WrapWithSpinner(fmt.Sprintf("Waiting 60s until %v interface is up, connected to peers and working", createdInterface.Created), func() error {
				deadline := time.Now().Add(time.Second * 60)
				for time.Now().Before(deadline) {
					myInterface, err := gokeenrestapi.Interface.GetInterfaceViaRciShowInterfaces(createdInterface.Created)
					if err != nil {
						return err
					}
					if myInterface.Connected == "yes" && myInterface.Link == "up" && myInterface.State == "up" {
						return nil
					}
					time.Sleep(time.Millisecond * 500)
				}
				return fmt.Errorf("looks like interface %v is no up. Please check The keenetic web-interface", createdInterface.Created)
			})
		}
		return err
	}
	return cmd
}
