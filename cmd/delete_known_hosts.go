package cmd

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/cobra"
)

func newDeleteKnownHostsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     CmdDeleteKnownHosts,
		Aliases: AliasesDeleteKnownHosts,
		Short:   "Delete known hosts by name or MAC using regex pattern",
	}

	var namePattern, macPattern string
	var force bool
	cmd.Flags().StringVar(&namePattern, "name-pattern", "", "Regex pattern to match host names for deletion")
	cmd.Flags().StringVar(&macPattern, "mac-pattern", "", "Regex pattern to match host MAC addresses for deletion")
	cmd.Flags().BoolVar(&force, "force", false, "Delete without confirmation")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if (namePattern == "") == (macPattern == "") {
			return fmt.Errorf("exactly one of --name-pattern or --mac-pattern must be specified")
		}

		var regex *regexp.Regexp
		var err error
		if namePattern != "" {
			regex, err = regexp.Compile(namePattern)
		} else {
			regex, err = regexp.Compile(macPattern)
		}
		if err != nil {
			return err
		}

		hotspot, err := gokeenrestapi.Ip.GetAllHotspots()
		if err != nil {
			return err
		}

		var hostMacsToDelete []string
		for _, host := range hotspot.Host {
			var matched bool
			if namePattern != "" {
				matched = regex.MatchString(host.Name)
			} else {
				matched = regex.MatchString(host.Mac)
			}
			if matched {
				gokeenlog.InfoSubStepf("Matching host: %v (MAC: %v)", color.CyanString(host.Name), color.BlueString(host.Mac))
				hostMacsToDelete = append(hostMacsToDelete, host.Mac)
			}
		}

		if len(hostMacsToDelete) == 0 {
			gokeenlog.Info("No hosts found matching the pattern, no need to delete")
			return nil
		}

		if !force {
			confirmed, err := confirmAction(fmt.Sprintf("\nFound %d host(s) to delete. Do you want to continue?", len(hostMacsToDelete)))
			if err != nil {
				return err
			}
			if !confirmed {
				gokeenlog.Info("Deletion cancelled")
				return nil
			}
		}

		return gokeenrestapi.Ip.DeleteKnownHosts(hostMacsToDelete)
	}
	return cmd
}
