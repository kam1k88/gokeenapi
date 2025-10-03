package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/gokeenlog"
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/spf13/cobra"
)

func newDeleteRoutesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete-routes",
		Aliases: []string{"deleteroutes", "dr"},
		Short:   "Delete static routes in Keenetic router",
	}

	var interfaceId string
	var force bool
	cmd.Flags().StringVar(&interfaceId, "interface-id", "", "Keenetic interface ID to delete static routes on, optional")
	cmd.Flags().BoolVar(&force, "force", false, "Delete without confirmation")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var interfaces []string
		if interfaceId != "" {
			interfaces = append(interfaces, interfaceId)
		} else {
			for _, routeSetting := range config.Cfg.Routes {
				interfaces = append(interfaces, routeSetting.InterfaceID)
			}
		}

		type interfaceRoutes struct {
			interfaceId string
			routes      []gokeenrestapimodels.RciIpRoute
		}
		var allRoutesToDelete []interfaceRoutes
		var totalRoutes int

		for _, ifaceId := range interfaces {
			err := gokeenrestapi.Checks.CheckInterfaceId(ifaceId)
			if err != nil {
				return err
			}
			err = gokeenrestapi.Checks.CheckInterfaceExists(ifaceId)
			if err != nil {
				return err
			}
			routes, err := gokeenrestapi.Ip.GetAllUserRoutesRciIpRoute(ifaceId)
			if err != nil {
				return err
			}

			if len(routes) > 0 {
				totalRoutes += len(routes)
				allRoutesToDelete = append(allRoutesToDelete, interfaceRoutes{ifaceId, routes})
			}
		}

		if totalRoutes == 0 {
			gokeenlog.Info("No routes found to delete")
			return nil
		}

		for _, routeInfo := range allRoutesToDelete {
			for _, route := range routeInfo.routes {
				msg := ""
				if route.Host != "" {
					msg = color.CyanString(route.Host)
				} else {
					msg = color.CyanString(route.Network) + "/" + color.BlueString(route.Mask)
				}
				gokeenlog.InfoSubStepf("Route to delete: %v via %v",
					msg,
					color.YellowString(route.Interface))
			}
		}

		if !force {
			confirmed, err := confirmAction(fmt.Sprintf("\nFound %v total route(s) to delete. Do you want to continue?", color.CyanString("%v", totalRoutes)))
			if err != nil {
				return err
			}
			if !confirmed {
				gokeenlog.Info("Deletion cancelled")
				return nil
			}
		}

		for _, item := range allRoutesToDelete {
			err := gokeenrestapi.Ip.DeleteRoutes(item.routes, item.interfaceId)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return cmd
}
