package cli

import (
	"errors"

	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic"
	"github.com/spf13/cobra"
)

func newUpdateAwgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     CmdUpdateAwg,
		Aliases: AliasesUpdateAwg,
		Short:   "Update existing WireGuard VPN configuration",
		Long: `Update an existing WireGuard (AWG) connection with new configuration from a .conf file.

This command updates the configuration of an existing WireGuard interface using
a standard WireGuard configuration file. It's useful for updating connection
parameters like endpoints, allowed IPs, or other peer settings.

The .conf file should contain standard WireGuard configuration:
  [Interface]
  PrivateKey = ...
  Address = ...
  
  [Peer]
  PublicKey = ...
  Endpoint = ...
  AllowedIPs = ...

Use cases:
- Update peer endpoint when server IP changes
- Modify allowed IPs for routing changes
- Update connection parameters without recreating interface
- Apply new WireGuard configuration to existing connection

Examples:
  # Update existing WireGuard interface
  gokeenapi update-awg --config config.yaml --conf-file /path/to/updated.conf --interface-id Wireguard0

  # Update with new server configuration
  gokeenapi update-awg --config config.yaml --conf-file /path/to/new-server.conf --interface-id MyVPN

Note: Use 'show-interfaces --type Wireguard' to list existing WireGuard interface IDs.`,
	}
	var confFile, interfaceId string
	cmd.Flags().StringVar(&confFile, "conf-file", "",
		`Path to WireGuard configuration file (.conf) with updated settings.
Must contain valid [Interface] and [Peer] sections.
This flag is required.`)
	cmd.Flags().StringVar(&interfaceId, "interface-id", "",
		`ID of the existing WireGuard interface to update.
Use 'show-interfaces --type Wireguard' to list available interfaces.
This flag is required.`)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if confFile == "" {
			return errors.New("--conf-file flag is required")
		}
		if interfaceId == "" {
			return errors.New("--interface-id flag is required")
		}
		return keenetic.AwgConf.ConfigureOrUpdateInterface(confFile, interfaceId)
	}
	return cmd
}
