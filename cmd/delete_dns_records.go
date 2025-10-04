package cmd

import (
	"fmt"
	"slices"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/internal/gokeenspinner"
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/spf13/cobra"
)

func newDeleteDnsRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     CmdDeleteDnsRecords,
		Aliases: AliasesDeleteDnsRecords,
		Short:   "Delete static dns records in Keenetic router",
	}

	var force bool
	cmd.Flags().BoolVar(&force, "force", false, "Delete without confirmation")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		runningConfig, err := gokeenrestapi.Common.ShowRunningConfig()
		if err != nil {
			return err
		}
		var parseC []gokeenrestapimodels.ParseRequest
		for _, addDnsRecordSetting := range config.Cfg.DNS.Records {
			for _, ip := range addDnsRecordSetting.IP {
				c := fmt.Sprintf("ip host %v %v", addDnsRecordSetting.Domain, ip)
				if !slices.Contains(runningConfig.Message, c) {
					continue
				}
				gokeenlog.InfoSubStepf("DNS record to delete: %v -> %v",
					color.CyanString(addDnsRecordSetting.Domain),
					color.BlueString(ip))
				c = fmt.Sprintf("no %v", c)
				parseC = append(parseC, gokeenrestapimodels.ParseRequest{Parse: c})
			}
		}
		if len(parseC) == 0 {
			gokeenlog.Info("No DNS records found to delete")
			return nil
		}

		if !force {
			confirmed, err := confirmAction(fmt.Sprintf("\nFound %v DNS record(s) to delete. Do you want to continue?", len(parseC)))
			if err != nil {
				return err
			}
			if !confirmed {
				gokeenlog.Info("Deletion cancelled")
				return nil
			}
		}
		err = gokeenspinner.WrapWithSpinner(fmt.Sprintf("Deleting %v DNS records", color.CyanString("%v", len(parseC))), func() error {
			parseC = append(parseC, gokeenrestapi.Common.SaveConfigParseRequest())
			result, err := gokeenrestapi.Common.ExecutePostParse(parseC...)
			if err != nil {
				return err
			}
			gokeenlog.PrintParseResponse(result)
			return err
		})
		return err
	}
	return cmd
}
