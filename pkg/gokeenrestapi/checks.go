package gokeenrestapi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"gopkg.in/ini.v1"
)

var (
	// Checks provides validation functions for router interface and configuration checks
	Checks = checks{}
)

type checks struct {
}

// CheckInterfaceId validates that the provided interface ID is not empty
func (*checks) CheckInterfaceId(interfaceId string) error {
	if interfaceId == "" {
		return errors.New("please specify a keenetic interface id via flag/field/variable")
	}
	return nil
}

// CheckInterfaceExists verifies that the specified interface ID exists on the router
func (*checks) CheckInterfaceExists(interfaceId string) error {
	interfaces, err := Interface.GetInterfacesViaRciShowInterfaces(false)
	if err != nil {
		return err
	}
	interfaceFound := false
	for _, interfaceDetails := range interfaces {
		if interfaceDetails.Id == interfaceId {
			interfaceFound = true
			break
		}
	}
	if !interfaceFound {
		return fmt.Errorf("keenetic router doesn't have interface with id '%v'. Verify that you specified correct ID", interfaceId)
	}
	return nil
}

// CheckComponentInstalled checks if a specific component is installed on the router
// Returns the installation status string if found, empty string if not installed or not found
func (*checks) CheckComponentInstalled(componentName string) (installed string, err error) {
	var components gokeenrestapimodels.RciComponentsList
	components, err = Components.GetAllComponents()
	if err != nil {
		return
	}
	for k, v := range components.Component {
		if strings.EqualFold(k, componentName) {
			installed = v.Installed
			return
		}
	}
	return
}

// CheckAWGInterfaceExistsFromConfFile checks if a WireGuard connection from the config file already exists
// Returns an error if a matching WireGuard interface is found with the same configuration
func (*checks) CheckAWGInterfaceExistsFromConfFile(confPath string) error {
	interfaces, err := Interface.GetInterfacesViaRciShowInterfaces(true, InterfaceTypeWireguard)
	if err != nil {
		return err
	}
	var interfacesIds []string
	for _, interfacesDetails := range interfaces {
		interfacesIds = append(interfacesIds, interfacesDetails.Id)
	}
	scInterfaces, err := Interface.GetInterfacesViaRciShowScInterfaces(interfacesIds...)
	if err != nil {
		return err
	}
	conf, err := ini.Load(confPath)
	if err != nil {
		return err
	}
	peerSection, err := conf.GetSection("Peer")
	if err != nil {
		return err
	}
	interfaceSection, err := conf.GetSection("Interface")
	if err != nil {
		return err
	}
	address, err := interfaceSection.GetKey("Address")
	if err != nil {
		return err
	}
	endpoint, err := peerSection.GetKey("Endpoint")
	if err != nil {
		return err
	}
	publicKey, err := peerSection.GetKey("PublicKey")
	if err != nil {
		return err
	}
	foundWgId := ""
	addressKey := address.String()
	addressKeySplit := strings.Split(addressKey, "/")
	publicKeyKey := publicKey.String()
	endpointKey := endpoint.String()
	for id, interfaceDetails := range scInterfaces {
		for _, peer := range interfaceDetails.Wireguard.Peer {
			if strings.EqualFold(publicKeyKey, peer.Key) && strings.EqualFold(endpointKey, peer.Endpoint.Address) &&
				strings.EqualFold(interfaceDetails.IP.Address.Address, addressKeySplit[0]) {
				foundWgId = id
				break
			}
		}
	}
	if foundWgId != "" {
		return fmt.Errorf("wireguard connection from the config file already exists with the following id: %v\n\nUse configure-awg command instead with --interface-id %v flag", color.CyanString(foundWgId), foundWgId)
	}
	return nil
}
