package gokeencache

import (
	"fmt"

	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/patrickmn/go-cache"
)

const (
	rciShowInterfaces = "rci_show_interfaces"
	runtimeConfig     = "runtime_config"
	rciShowIpRoute    = "rci_show_ip_route"
)

var (
	c = cache.New(cache.NoExpiration, cache.NoExpiration)
)

func UpdateRuntimeConfig(f func(runtime *config.Runtime)) {
	cfg := GetRuntimeConfig()
	f(cfg)
	c.Set(runtimeConfig, cfg, cache.NoExpiration)
}

func GetRuntimeConfig() *config.Runtime {
	cfg, ok := c.Get(runtimeConfig)
	if ok {
		return cfg.(*config.Runtime)
	}
	return &config.Runtime{}
}

func SetRciShowIpRoute(routes []gokeenrestapimodels.RciShowIpRoute, interfaceId string) {
	if interfaceId == "" {
		interfaceId = "all"
	}
	c.Set(fmt.Sprintf("%v-%v", rciShowIpRoute, interfaceId), routes, cache.NoExpiration)
}

func GetRciShowIpRoute(interfaceId string) []gokeenrestapimodels.RciShowIpRoute {
	if interfaceId == "" {
		interfaceId = "all"
	}
	v, ok := c.Get(fmt.Sprintf("%v-%v", rciShowIpRoute, interfaceId))
	if !ok {
		return nil
	}
	return v.([]gokeenrestapimodels.RciShowIpRoute)
}

func SetRciShowInterfaces(m map[string]gokeenrestapimodels.RciShowInterface) {
	c.Set(rciShowInterfaces, m, cache.NoExpiration)
}
func GetRciShowInterfaces() map[string]gokeenrestapimodels.RciShowInterface {
	v, ok := c.Get(rciShowInterfaces)
	if !ok {
		return nil
	}
	return v.(map[string]gokeenrestapimodels.RciShowInterface)
}
