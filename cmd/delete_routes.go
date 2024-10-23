package cmd

import (
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

	cmd.Flags().String("interface", "", "Keenetic interface ID to delete static routes on")
	_ = cmd.MarkFlagRequired("interface")
	_ = viper.BindPFlag(config.ViperKeeneticInterface, cmd.Flags().Lookup("interface"))
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		routes, err := keeneticapi.Route.GetAllUserRoutesRciIpRoute(viper.GetString(config.ViperKeeneticInterface))
		if err != nil {
			return err
		}
		return keeneticapi.Route.DeleteRoutes(routes)
	}
	return cmd
}
