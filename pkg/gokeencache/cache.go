package gokeencache

import (
	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"github.com/patrickmn/go-cache"
)

var (
	c = cache.New(cache.NoExpiration, cache.NoExpiration)
)

func SetVersion(v gokeenrestapimodels.Version) {
	c.Set("version", v, cache.NoExpiration)
}

func GetVersion() gokeenrestapimodels.Version {
	v, _ := c.Get("version")
	return v.(gokeenrestapimodels.Version)
}
