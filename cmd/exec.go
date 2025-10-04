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

// execCommand executes a command on the Keenetic router and returns the results
func execCommand(args []string) ([]gokeenrestapimodels.ParseResponse, error) {
	cmdToExecute := strings.Join(args, " ")
	parseC := gokeenrestapimodels.ParseRequest{Parse: cmdToExecute}
	return gokeenrestapi.Common.ExecutePostParse(parseC)
}

// printExecResults prints the execution results to stdout
func printExecResults(results []gokeenrestapimodels.ParseResponse) {
	gokeenlog.Info("Result:")
	for _, r := range results {
		if r.Parse.DynamicData != "" {
			fmt.Println(r.Parse.DynamicData)
		}
	}
}

func newExecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     CmdExec,
		Aliases: AliasesExec,
		Short:   "Execute any custom command in Keenetic router",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cmdToExecute := strings.Join(args, " ")
		return gokeenspinner.WrapWithSpinner(fmt.Sprintf("Executing %v command", color.CyanString(cmdToExecute)), func() error {
			result, err := execCommand(args)
			if err != nil {
				return err
			}
			printExecResults(result)
			return nil
		})
	}
	return cmd
}
