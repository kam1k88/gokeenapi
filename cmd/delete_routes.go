package cmd

import (
	"fmt"
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/pkg/keeneticapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newDeleteRoutesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete-routes",
		Aliases: []string{"deleteroutes"},
		Short:   "Delete static routes in Keenetic router",
	}

	cmd.Flags().String("interface-id", "", "Keenetic interface ID to delete static routes on")
	_ = cmd.MarkFlagRequired("interface-id")

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag(config.ViperKeeneticInterfaceId, cmd.Flags().Lookup("interface-id"))
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		interfaces, err := keeneticapi.Interface.GetInterfacesViaRciShowInterfaces()
		if err != nil {
			return err
		}
		interfaceFound := false
		for _, interfaceDetails := range interfaces {
			if interfaceDetails.Id == viper.GetString(config.ViperKeeneticInterfaceId) {
				interfaceFound = true
				break
			}
		}
		if !interfaceFound {
			return fmt.Errorf("keenetic router doesn't have interface with id '%v'", viper.GetString(config.ViperKeeneticInterfaceId))
		}
		routes, err := keeneticapi.Route.GetAllUserRoutesRciIpRoute(viper.GetString(config.ViperKeeneticInterfaceId))
		if err != nil {
			return err
		}
		return keeneticapi.Route.DeleteRoutes(routes)
	}
	return cmd
}
