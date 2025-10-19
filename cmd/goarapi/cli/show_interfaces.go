package cli

import (
	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic"
	"github.com/spf13/cobra"
)

func newShowInterfacesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     CmdShowInterfaces,
		Aliases: AliasesShowInterfaces,
		Short:   "List available network interfaces on your router",
		Long: `Display detailed information about network interfaces available on your Keenetic router.

This command helps you discover interface IDs that can be used with other commands like add-routes 
and delete-routes. It shows interface names, types, status, and other relevant details.

Examples:
  # Show all interfaces
  gokeenapi show-interfaces --config config.yaml

  # Show only WireGuard interfaces  
  gokeenapi show-interfaces --config config.yaml --type Wireguard

  # Show multiple interface types
  gokeenapi show-interfaces --config config.yaml --type Wireguard --type Ethernet`,
	}

	var interfaceType []string
	cmd.Flags().StringSliceVar(&interfaceType, "type", []string{},
		`Filter interfaces by type (e.g., Wireguard, Ethernet, Bridge).
Can be specified multiple times to show multiple types.
If not specified, shows all interface types.`)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		interfaces, err := keenetic.Interface.GetInterfacesViaRciShowInterfaces(true, interfaceType...)
		if err != nil {
			return err
		}
		keenetic.Interface.PrintInfoAboutInterfaces(interfaces)
		return nil
	}
	return cmd
}
