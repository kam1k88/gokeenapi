package cmd

import (
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/cobra"
)

func newDeleteRoutesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete-routes",
		Aliases: []string{"deleteroutes", "dr"},
		Short:   "Delete static routes in Keenetic router",
	}

	var interfaceId string
	cmd.Flags().StringVar(&interfaceId, "interface-id", "", "Keenetic interface ID to delete static routes on, optional")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var interfaces []string
		if interfaceId != "" {
			interfaces = append(interfaces, interfaceId)
		} else {
			for _, routeSetting := range config.Cfg.Routes {
				interfaces = append(interfaces, routeSetting.InterfaceID)
			}
		}
		for _, ifaceId := range interfaces {
			err := gokeenrestapi.Checks.CheckInterfaceId(ifaceId)
			if err != nil {
				return err
			}
			err = gokeenrestapi.Checks.CheckInterfaceExists(ifaceId)
			if err != nil {
				return err
			}
			routes, err := gokeenrestapi.Ip.GetAllUserRoutesRciIpRoute(ifaceId)
			if err != nil {
				return err
			}
			err = gokeenrestapi.Ip.DeleteRoutes(routes, ifaceId)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return cmd
}
