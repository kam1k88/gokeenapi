package cmd

import (
	"errors"

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
	cmd.Flags().StringVar(&interfaceId, "interface-id", "", "Keenetic interface ID to delete static routes on")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if interfaceId == "" {
			return errors.New("--interface-id flag is required")
		}
		err := gokeenrestapi.Checks.CheckInterfaceId(interfaceId)
		if err != nil {
			return err
		}
		err = gokeenrestapi.Checks.CheckInterfaceExists(interfaceId)
		if err != nil {
			return err
		}
		routes, err := gokeenrestapi.Ip.GetAllUserRoutesRciIpRoute(interfaceId)
		if err != nil {
			return err
		}
		return gokeenrestapi.Ip.DeleteRoutes(routes, interfaceId)
	}
	return cmd
}
