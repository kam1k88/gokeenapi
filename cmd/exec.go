package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/internal/gokeenspinner"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/spf13/cobra"
)

func newExecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exec",
		Aliases: []string{"e"},
		Short:   "Execute any custom command in Keenetic router",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cmdToExecute := strings.Join(args, " ")
		parseC := gokeenrestapimodels.ParseRequest{Parse: cmdToExecute}
		return gokeenspinner.WrapWithSpinner(fmt.Sprintf("Executing %v command", color.CyanString(cmdToExecute)), func() error {
			result, err := gokeenrestapi.Common.ExecutePostParse(parseC)
			if err != nil {
				return err
			}
			gokeenlog.Info("Result:")
			for _, r := range result {
				if r.Parse.DynamicData != "" {
					fmt.Println(r.Parse.DynamicData)
				}
			}
			return nil
		})
	}
	return cmd
}
