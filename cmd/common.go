package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapi"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
	"gopkg.in/ini.v1"
)

func checkRequiredFields() error {
	var mErr error
	if viper.GetString(config.ViperKeeneticUrl) == "" {
		mErr = multierr.Append(mErr, errors.New("please specify a keenetic url via flag/field/variable"))
	}
	if viper.GetString(config.ViperKeeneticLogin) == "" {
		mErr = multierr.Append(mErr, errors.New("please specify a keenetic login via flag/field/variable"))
	}
	if viper.GetString(config.ViperKeeneticPassword) == "" {
		mErr = multierr.Append(mErr, errors.New("please specify a keenetic password via flag/field/variable"))
	}

	return mErr
}

func checkInterfaceId(interfaceId string) error {
	if interfaceId == "" {
		return errors.New("please specify a keenetic interface id via flag/field/variable")
	}
	return nil
}

func checkInterfaceExists(interfaceId string) error {
	interfaces, err := gokeenrestapi.Interface.GetInterfacesViaRciShowInterfaces()
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

func checkAWGInterfaceExistsFromConfFile(confPath string) error {
	interfaces, err := gokeenrestapi.Interface.GetInterfacesViaRciShowInterfaces("Wireguard")
	if err != nil {
		return err
	}
	var interfacesIds []string
	for _, interfacesDetails := range interfaces {
		interfacesIds = append(interfacesIds, interfacesDetails.Id)
	}
	scInterfaces, err := gokeenrestapi.Interface.GetInterfacesViaRciShowScInterfaces(interfacesIds...)
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

func RestoreCursor() {
	if !(len(os.Getenv("WT_SESSION")) > 0 && runtime.GOOS == "windows") {
		// make sure to restore cursor in all cases
		_, _ = fmt.Fprint(os.Stdout, "\033[?25h")
	}
}
