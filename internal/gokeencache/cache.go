package gokeencache

import (
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/patrickmn/go-cache"
)

const (
	rciShowInterfaces = "rci_show_interfaces"
	version           = "version"
)

var (
	c = cache.New(cache.NoExpiration, cache.NoExpiration)
)

func SetVersion(v gokeenrestapimodels.Version) {
	c.Set(version, v, cache.NoExpiration)
}

func GetVersion() gokeenrestapimodels.Version {
	v, _ := c.Get(version)
	return v.(gokeenrestapimodels.Version)
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
