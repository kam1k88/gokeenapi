package cmd

import (
	"fmt"
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/pkg/keeneticapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

func newAddRoutesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-routes",
		Aliases: []string{"addroutes"},
		Short:   "Add static routes in Keenetic router",
	}

	var batFiles []string
	var batUrls []string
	cmd.Flags().String("interface-id", "", "Keenetic interface ID to update routes on")
	cmd.Flags().StringSliceVar(&batFiles, "bat-file", []string{}, "Path to a bat file to add routes from. Can be specified multiple times")
	cmd.Flags().StringSliceVar(&batUrls, "bat-url", []string{}, "URL to a bat file to add routes from. Can be specified multiple times")
	_ = cmd.MarkFlagRequired("interface")

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag(config.ViperKeeneticInterfaceId, cmd.Flags().Lookup("interface-id"))
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(batFiles) == 0 && len(batUrls) == 0 {
			return fmt.Errorf("at least one of --bat-file or --bat-url must be set")
		}
		for _, file := range batFiles {
			absFilePath, err := filepath.Abs(file)
			if err != nil {
				return err
			}
			err = keeneticapi.Route.AddRoutesFromBatFile(absFilePath)
			if err != nil {
				return err
			}
		}
		for _, url := range batUrls {
			err := keeneticapi.Route.AddRoutesFromBatUrl(url)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return cmd
}
