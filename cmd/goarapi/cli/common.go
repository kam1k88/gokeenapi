package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic"
	"github.com/kam1k88/gokeenapi/pkg/config"
	"github.com/kam1k88/gokeenapi/pkg/goarapi"
	"github.com/kam1k88/gokeenapi/pkg/registry"
	"go.uber.org/multierr"
)

var (
	routerRegistry = registry.New()
	apiFacade      *goarapi.AnyRouterAPI
)

func init() {
	routerRegistry.RegisterBackend("keenetic", func(name string, cfg *config.GokeenapiConfig) (goarapi.RouterAPI, error) {
		return keenetic.NewRouter(name, cfg), nil
	})
}

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

func prepareAPI(ctx context.Context) error {
	apiFacade = goarapi.New()
	router, err := routerRegistry.Create("keenetic", "keenetic", &config.Cfg)
	if err != nil {
		return err
	}
	apiFacade.Register("keenetic", router)
	return apiFacade.AuthenticateAll(ctx)
}

// API exposes initialized facade for other commands.
func API() *goarapi.AnyRouterAPI {
	return apiFacade
}

func confirmAction(message string) (bool, error) {
	fmt.Printf("%s (y/N): ", message)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}
