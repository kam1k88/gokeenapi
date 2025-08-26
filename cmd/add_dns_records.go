package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/internal/gokeenspinner"
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/spf13/cobra"
)

func newAddDnsRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-dns-records",
		Aliases: []string{"adddnsrecords", "adr"},
		Short:   "Add static dns records in Keenetic router",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		//runningConfig, err := gokeenrestapi.Common.ShowRunningConfig()
		//if err != nil {
		//	return err
		//}
		var parseC []gokeenrestapimodels.ParseRequest
		for _, addDnsRecordSetting := range config.Cfg.DNS.Records {
			for _, ip := range addDnsRecordSetting.IP {
				c := fmt.Sprintf("ip host %v %v", addDnsRecordSetting.Domain, ip)
				parseC = append(parseC, gokeenrestapimodels.ParseRequest{Parse: c})
			}
		}
		if len(parseC) == 0 {
			return nil
		}
		err := gokeenspinner.WrapWithSpinner(fmt.Sprintf("Adding %v DNS records", color.CyanString("%v", len(parseC))), func() error {
			parseC = append(parseC, gokeenrestapimodels.ParseRequest{Parse: "system configuration save"})
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
