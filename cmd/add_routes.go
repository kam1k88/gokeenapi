package cmd

import (
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/pkg/keeneticapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAddRoutesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-routes",
		Aliases: []string{"addroutes"},
		Short:   "Add static routes in Keenetic router",
	}

	var batFile string
	cmd.Flags().String("interface", "", "Keenetic interface ID to update routes on")
	cmd.Flags().StringVar(&batFile, "bat-file", "", "Path to a bat file to add routes from")
	_ = cmd.MarkFlagRequired("interface")
	_ = viper.BindPFlag(config.ViperKeeneticInterface, cmd.Flags().Lookup("interface"))
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if batFile != "" {
			err := keeneticapi.Route.AddRoutesFromBatFile(batFile)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return cmd
}
