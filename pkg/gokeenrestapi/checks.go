package gokeenrestapi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/ini.v1"
)

var (
	Checks = checks{}
)

type checks struct {
}

func (*checks) CheckInterfaceId(interfaceId string) error {
	if interfaceId == "" {
		return errors.New("please specify a keenetic interface id via flag/field/variable")
	}
	return nil
}

func (*checks) CheckInterfaceExists(interfaceId string) error {
	interfaces, err := Interface.GetInterfacesViaRciShowInterfaces(true)
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

func (*checks) CheckAWGInterfaceExistsFromConfFile(confPath string) error {
	interfaces, err := Interface.GetInterfacesViaRciShowInterfaces(true, "Wireguard")
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
