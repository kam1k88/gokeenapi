package cmd

import (
	"errors"

	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/cobra"
)

func newConfigureAwgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "configure-awg",
		Aliases: []string{"configureawg", "cawg"},
		Short:   "Configure AWG connection to add(update) ASC parameters in it in Keenetic router",
	}
	var confFile, interfaceId string
	cmd.Flags().StringVar(&confFile, "conf-file", "", "Path to a conf file with AWG configuration")
	cmd.Flags().StringVar(&interfaceId, "interface-id", "", "ID of existing interface to update from the conf file")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if confFile == "" {
			return errors.New("--conf-file flag is required")
		}
		if interfaceId == "" {
			return errors.New("--interface-id flag is required")
		}
		return gokeenrestapi.AwgConf.ConfigureOrUpdateInterface(confFile, interfaceId)
	}
	return cmd
}
