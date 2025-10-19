package gokeencache

import (
	"github.com/kam1k88/gokeenapi/pkg/backends/keenetic/models"
	"github.com/kam1k88/gokeenapi/pkg/config"
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

func SetRciShowIpRoute(routes []models.RciShowIpRoute) {
	c.Set(rciShowIpRoute, routes, cache.NoExpiration)
}

func GetRciShowIpRoute() []models.RciShowIpRoute {
	v, ok := c.Get(rciShowIpRoute)
	if !ok {
		return nil
	}
	return v.([]models.RciShowIpRoute)
}

func SetRciShowInterfaces(m map[string]models.RciShowInterface) {
	c.Set(rciShowInterfaces, m, cache.NoExpiration)
}
func GetRciShowInterfaces() map[string]models.RciShowInterface {
	v, ok := c.Get(rciShowInterfaces)
	if !ok {
		return nil
	}
	return v.(map[string]models.RciShowInterface)
}
