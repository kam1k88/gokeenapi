package cmd

import (
	"path/filepath"

	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/cobra"
)

func newAddRoutesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     CmdAddRoutes,
		Aliases: AliasesAddRoutes,
		Short:   "Add static routes in Keenetic router",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		for _, addRouteSettings := range config.Cfg.Routes {
			err := gokeenrestapi.Checks.CheckInterfaceId(addRouteSettings.InterfaceID)
			if err != nil {
				return err
			}
			err = gokeenrestapi.Checks.CheckInterfaceExists(addRouteSettings.InterfaceID)
			if err != nil {
				return err
			}
			for _, file := range addRouteSettings.BatFile {
				absFilePath, err := filepath.Abs(file)
				if err != nil {
					return err
				}
				err = gokeenrestapi.Ip.AddRoutesFromBatFile(absFilePath, addRouteSettings.InterfaceID)
				if err != nil {
					return err
				}
			}
			for _, url := range addRouteSettings.BatURL {
				err := gokeenrestapi.Ip.AddRoutesFromBatUrl(url, addRouteSettings.InterfaceID)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	return cmd
}
