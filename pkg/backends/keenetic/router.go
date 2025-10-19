package keenetic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/fatih/color"

	"github.com/kam1k88/gokeenapi/internal/gokeencache"
	"github.com/kam1k88/gokeenapi/internal/gokeenlog"
	"github.com/kam1k88/gokeenapi/internal/gokeenspinner"
	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic/models"
	"github.com/kam1k88/gokeenapi/pkg/config"
	"github.com/kam1k88/gokeenapi/pkg/goarapi"
)

// Router implements goarapi.RouterAPI for Keenetic devices.
type Router struct {
	name string
	cfg  *config.GokeenapiConfig
}

// NewRouter creates new router instance bound to provided logical name and configuration snapshot.
func NewRouter(name string, cfg *config.GokeenapiConfig) *Router {
	return &Router{name: name, cfg: cfg}
}

// Name returns logical router name.
func (r *Router) Name() string {
	return r.name
}

// Authenticate ensures configuration is applied globally and performs login via Keenetic common facade.
func (r *Router) Authenticate(ctx context.Context) error {
	if r.cfg != nil {
		config.Cfg = *r.cfg
	}
	return Common.Auth()
}

// DeviceInfo returns router metadata obtained from Keenetic API.
func (r *Router) DeviceInfo(ctx context.Context) (goarapi.DeviceInfo, error) {
	version, err := Common.Version()
	if err != nil {
		return goarapi.DeviceInfo{}, err
	}
	return goarapi.DeviceInfo{
		Name:     r.name,
		Model:    version.Model,
		Version:  version.Title,
		Firmware: version.Release,
	}, nil
}

// ListRoutes queries Keenetic for configured static routes.
func (r *Router) ListRoutes(ctx context.Context) ([]goarapi.Route, error) {
	if r.cfg != nil {
		config.Cfg = *r.cfg
	}
	body, err := Common.ExecuteGetSubPath("/rci/ip/route")
	if err != nil {
		return nil, err
	}
	var items []models.RciIpRoute
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	routes := make([]goarapi.Route, 0, len(items))
	for _, item := range items {
		id := buildRouteID(item)
		network := item.Network
		mask := item.Mask
		host := item.Host
		gateway := "auto"
		if !item.Auto {
			gateway = "manual"
		}
		routes = append(routes, goarapi.Route{
			ID:        id,
			Network:   network,
			Mask:      mask,
			Host:      host,
			Gateway:   gateway,
			Interface: item.Interface,
		})
	}
	return routes, nil
}

// AddRoute provisions new static route via Keenetic parse API.
func (r *Router) AddRoute(ctx context.Context, route goarapi.Route) error {
	if r.cfg != nil {
		config.Cfg = *r.cfg
	}
	if route.Interface == "" {
		return fmt.Errorf("interface is required")
	}
	if route.Network == "" || route.Mask == "" {
		return fmt.Errorf("network and mask are required")
	}
	gateway := route.Gateway
	if gateway == "" {
		gateway = "auto"
	}
	command := fmt.Sprintf("ip route %s %s %s %s", route.Network, route.Mask, route.Interface, gateway)
	parse := []models.ParseRequest{{Parse: command}}
	parse = Common.EnsureSaveConfigAtEnd(parse)
	gokeencache.SetRciShowIpRoute(nil)
	return gokeenspinner.WrapWithSpinner(fmt.Sprintf("Adding static route %s/%s", color.CyanString(route.Network), route.Mask), func() error {
		_, err := Common.ExecutePostParse(parse...)
		return err
	})
}

// DeleteRoute removes a static route identified by key (Route.ID).
func (r *Router) DeleteRoute(ctx context.Context, key string) error {
	if r.cfg != nil {
		config.Cfg = *r.cfg
	}
	decoded, err := url.PathUnescape(key)
	if err != nil {
		decoded = key
	}
	iface, targetType, network, mask, err := parseRouteID(decoded)
	if err != nil {
		return err
	}
	var spec string
	switch targetType {
	case "host":
		spec = network
	case "net":
		if mask == "" {
			return fmt.Errorf("mask is required for network route")
		}
		spec = fmt.Sprintf("%s %s", network, mask)
	default:
		return fmt.Errorf("unknown route type %s", targetType)
	}
	command := fmt.Sprintf("no ip route %s %s", spec, iface)
	parse := []models.ParseRequest{{Parse: command}}
	parse = Common.EnsureSaveConfigAtEnd(parse)
	gokeencache.SetRciShowIpRoute(nil)
	return gokeenspinner.WrapWithSpinner(fmt.Sprintf("Deleting static route %s", color.CyanString(spec)), func() error {
		_, err := Common.ExecutePostParse(parse...)
		if err != nil {
			return err
		}
		gokeenlog.InfoSubStepf("Route %s removed from %s", spec, iface)
		return nil
	})
}

func buildRouteID(route models.RciIpRoute) string {
	if route.Host != "" {
		return fmt.Sprintf("%s|host|%s", route.Interface, route.Host)
	}
	return fmt.Sprintf("%s|net|%s|%s", route.Interface, route.Network, route.Mask)
}

func parseRouteID(id string) (iface, routeType, network, mask string, err error) {
	parts := strings.Split(id, "|")
	if len(parts) < 3 {
		return "", "", "", "", fmt.Errorf("invalid route identifier: %s", id)
	}
	iface = parts[0]
	routeType = parts[1]
	switch routeType {
	case "host":
		network = strings.Join(parts[2:], "|")
	case "net":
		if len(parts) < 4 {
			return "", "", "", "", fmt.Errorf("invalid network route identifier: %s", id)
		}
		network = parts[2]
		mask = strings.Join(parts[3:], "|")
	default:
		err = fmt.Errorf("unknown route type %s", routeType)
	}
	return
}
