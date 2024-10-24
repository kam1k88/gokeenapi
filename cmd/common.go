package cmd

import (
	"errors"
	"github.com/noksa/gokeenapi/internal/config"
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
