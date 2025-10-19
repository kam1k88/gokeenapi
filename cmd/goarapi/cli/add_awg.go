package cli

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/kam1k88/gokeenapi/internal/gokeenlog"
	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic"
	"github.com/spf13/cobra"
)

func newAddAwgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     CmdAddAwg,
		Aliases: AliasesAddAwg,
		Short:   "Set up a new WireGuard VPN connection",
		Long: `Add a new WireGuard (AWG) VPN connection to your Keenetic router from a .conf file.

This command creates and configures a new WireGuard interface using a standard 
WireGuard configuration file. The interface is automatically configured, enabled,
and brought up after creation.

The .conf file should contain standard WireGuard configuration:
  [Interface]
  PrivateKey = ...
  Address = ...
  
  [Peer]
  PublicKey = ...
  Endpoint = ...
  AllowedIPs = ...

Process:
1. Validates the configuration file
2. Creates a new WireGuard interface
3. Applies the configuration
4. Enables global IP routing
5. Brings the interface up
6. Waits for interface to become active

Examples:
  # Add WireGuard connection with auto-generated name
  gokeenapi add-awg --config config.yaml --conf-file /path/to/wg.conf

  # Add WireGuard connection with custom name
  gokeenapi add-awg --config config.yaml --conf-file /path/to/wg.conf --name MyVPN`,
	}
	var name, confFile string
	cmd.Flags().StringVar(&confFile, "conf-file", "",
		`Path to WireGuard configuration file (.conf).
Must contain valid [Interface] and [Peer] sections.
This flag is required.`)
	cmd.Flags().StringVar(&name, "name", "",
		`Custom name for the new WireGuard interface.
If not specified, Keenetic will auto-generate a name (e.g., Wireguard0).
The name will be used as the interface ID for other commands.`)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if confFile == "" {
			return fmt.Errorf("conf-file flag is required")
		}
		confPath, err := filepath.Abs(confFile)
		if err != nil {
			return err
		}
		err = keenetic.Checks.CheckAWGInterfaceExistsFromConfFile(confPath)
		if err != nil {
			return err
		}
		gokeenlog.InfoSubStepf("Conf-file: %v", color.CyanString("%v", confPath))
		createdInterface, err := keenetic.AwgConf.AddInterface(confPath, name)
		if err != nil {
			return err
		}
		gokeenlog.InfoSubStepf("ID: %v", color.CyanString(createdInterface.Created))
		err = keenetic.AwgConf.ConfigureOrUpdateInterface(confPath, createdInterface.Created)
		if err != nil {
			return err
		}
		err = keenetic.Interface.SetGlobalIpInInterface(createdInterface.Created, true)
		if err != nil {
			return err
		}
		err = keenetic.Interface.UpInterface(createdInterface.Created)
		if err != nil {
			return err
		}
		err = keenetic.Interface.WaitUntilInterfaceIsUp(createdInterface.Created)
		return err
	}
	return cmd
}
