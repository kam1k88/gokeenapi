package cmd

import (
	"github.com/noksa/gokeenapi/pkg/keeneticapi"
	"github.com/spf13/cobra"
)

func newShowInterfacesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show-interfaces",
		Aliases: []string{"showinterfaces", "showifaces"},
		Short:   "Print common information about interfaces in Keenetic router",
	}

	var interfaceType []string
	cmd.Flags().StringSliceVar(&interfaceType, "type", []string{}, "Show information about interfaces with specified type. Can be specified multiple times")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		interfaces, err := keeneticapi.Interface.GetInterfacesViaRciShowInterfaces(interfaceType...)
		if err != nil {
			return err
		}
		keeneticapi.Interface.PrintInfoAboutInterfaces(interfaces)
		return nil
	}
	return cmd
}
