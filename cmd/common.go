package cmd

import (
	"errors"
	"fmt"
	"github.com/noksa/gokeenapi/internal/config"
	"github.com/noksa/gokeenapi/pkg/keeneticapi"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

func checkRequiredFields() error {
	var mErr error
	if viper.GetString(config.ViperKeeneticUrl) == "" {
		mErr = multierr.Append(mErr, errors.New("please specify a keenetic api flag/field/variable"))
	}
	if viper.GetString(config.ViperKeeneticLogin) == "" {
		mErr = multierr.Append(mErr, errors.New("please specify a keenetic login flag/field/variable"))
	}
	if viper.GetString(config.ViperKeeneticPassword) == "" {
		mErr = multierr.Append(mErr, errors.New("please specify a keenetic password flag/field/variable"))
	}

	return mErr
}

func checkInterfaceExists() error {
	interfaces, err := keeneticapi.Interface.GetInterfacesViaRciShowInterfaces()
	if err != nil {
		return err
	}
	interfaceFound := false
	for _, interfaceDetails := range interfaces {
		if interfaceDetails.Id == viper.GetString(config.ViperKeeneticInterfaceId) {
			interfaceFound = true
			break
		}
	}
	if !interfaceFound {
		return fmt.Errorf("keenetic router doesn't have interface with id '%v'. Verify that you specified correct ID", viper.GetString(config.ViperKeeneticInterfaceId))
	}
	return nil
}
