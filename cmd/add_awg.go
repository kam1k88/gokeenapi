package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/cobra"
)

func newAddAwgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     CmdAddAwg,
		Aliases: AliasesAddAwg,
		Short:   "Add AWG connection from conf file in Keenetic router",
	}
	var name, confFile string
	cmd.Flags().StringVar(&confFile, "conf-file", "", "Path to a conf file with AWG configuration")
	cmd.Flags().StringVar(&name, "name", "", "Name for new WG interface, optional")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if confFile == "" {
			return fmt.Errorf("conf-file flag is required")
		}
		confPath, err := filepath.Abs(confFile)
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
		gokeenlog.InfoSubStepf("ID: %v", color.CyanString(createdInterface.Created))
		err = gokeenrestapi.AwgConf.ConfigureOrUpdateInterface(confPath, createdInterface.Created)
		if err != nil {
			return err
		}
		err = gokeenrestapi.Interface.SetGlobalIpInInterface(createdInterface.Created, true)
		if err != nil {
			return err
		}
		err = gokeenrestapi.Interface.UpInterface(createdInterface.Created)
		if err != nil {
			return err
		}
		err = gokeenrestapi.Interface.WaitUntilInterfaceIsUp(createdInterface.Created)
		return err
	}
	return cmd
}
