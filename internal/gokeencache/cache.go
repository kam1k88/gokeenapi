package gokeencache

import (
	"github.com/noksa/gokeenapi/pkg/config"
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/patrickmn/go-cache"
)

const (
	rciShowInterfaces = "rci_show_interfaces"
	runtimeConfig     = "runtime_config"
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
