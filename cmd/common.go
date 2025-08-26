package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/noksa/gokeenapi/pkg/config"
	"go.uber.org/multierr"
)

func checkRequiredFields() error {
	var mErr error
	if config.Cfg.Keenetic.URL == "" {
		mErr = multierr.Append(mErr, errors.New("please specify a keenetic url via keenetic.url field in yaml config"))
	}
	if config.Cfg.Keenetic.Login == "" {
		mErr = multierr.Append(mErr, errors.New("please specify a keenetic login via keenetic.login field in yaml config"))
	}
	if config.Cfg.Keenetic.Password == "" {
		mErr = multierr.Append(mErr, errors.New("please specify a keenetic password via keenetic.password field in yaml config"))
	}

	return mErr
}

func RestoreCursor() {
	if !(len(os.Getenv("WT_SESSION")) > 0 && runtime.GOOS == "windows") {
		// make sure to restore cursor in all cases
		_, _ = fmt.Fprint(os.Stdout, "\033[?25h")
	}
}
