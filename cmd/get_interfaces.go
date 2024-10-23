package cmd

import (
	"github.com/noksa/gokeenapi/pkg/keeneticapi"
	"github.com/spf13/cobra"
)

func newGetInterfacesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get-interfaces",
		Aliases: []string{"getinterfaces", "getifaces"},
		Short:   "Get common information about interfaces in Keenetic router",
	}

	var interfaceType string
	cmd.Flags().StringVar(&interfaceType, "type", "", "Show information about interfaces with specified type")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, err := keeneticapi.Interface.RciShowInterfaces(interfaceType)
		return err
	}
	return cmd
}
