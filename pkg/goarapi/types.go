package goarapi

import "context"

// RouterAPI describes the minimal set of capabilities that every router backend
// must implement to be managed by the AnyRouterAPI facade.
type RouterAPI interface {
	// Name returns logical router name used for registration and lookups.
	Name() string
	// Authenticate prepares backend for use (establishes sessions, etc.).
	Authenticate(ctx context.Context) error
	// DeviceInfo returns router metadata like model and firmware version.
	DeviceInfo(ctx context.Context) (DeviceInfo, error)
	// ListRoutes returns a snapshot of configured routes on the router.
	ListRoutes(ctx context.Context) ([]Route, error)
	// AddRoute creates a new static route on the router.
	AddRoute(ctx context.Context, route Route) error
	// DeleteRoute removes an existing route, identified by its key (network/host).
	DeleteRoute(ctx context.Context, key string) error
}

// DeviceInfo provides basic metadata about a router instance.
type DeviceInfo struct {
	Name     string `json:"name"`
	Model    string `json:"model"`
	Version  string `json:"version"`
	Firmware string `json:"firmware,omitempty"`
}

// Route represents a static route entry returned by router backends.
type Route struct {
	ID        string `json:"id"`
	Network   string `json:"network,omitempty"`
	Mask      string `json:"mask,omitempty"`
	Host      string `json:"host,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Interface string `json:"interface"`
	Distance  int    `json:"distance,omitempty"`
}
