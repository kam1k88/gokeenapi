package goarapi

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// AnyRouterAPI acts as a facade that aggregates multiple router backends and
// exposes a unified interface for higher level components (CLI, REST server).
type AnyRouterAPI struct {
	mu      sync.RWMutex
	routers map[string]RouterAPI
}

// New creates an empty AnyRouterAPI instance.
func New() *AnyRouterAPI {
	return &AnyRouterAPI{routers: make(map[string]RouterAPI)}
}

// Register adds a router backend under the provided name. Existing entries are
// replaced which allows reloading configuration at runtime.
func (a *AnyRouterAPI) Register(name string, router RouterAPI) {
	if router == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.routers[name] = router
}

// Router retrieves a router backend by name.
func (a *AnyRouterAPI) Router(name string) (RouterAPI, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	router, ok := a.routers[name]
	return router, ok
}

// Routers returns sorted list of registered router names.
func (a *AnyRouterAPI) Routers() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	names := make([]string, 0, len(a.routers))
	for name := range a.routers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// AuthenticateAll establishes sessions for all registered routers.
func (a *AnyRouterAPI) AuthenticateAll(ctx context.Context) error {
	a.mu.RLock()
	routers := make([]RouterAPI, 0, len(a.routers))
	for _, router := range a.routers {
		routers = append(routers, router)
	}
	a.mu.RUnlock()

	for _, router := range routers {
		if err := router.Authenticate(ctx); err != nil {
			return fmt.Errorf("authenticate router %s: %w", router.Name(), err)
		}
	}
	return nil
}

// DeviceInfo returns metadata for a specific router.
func (a *AnyRouterAPI) DeviceInfo(ctx context.Context, name string) (DeviceInfo, error) {
	router, ok := a.Router(name)
	if !ok {
		return DeviceInfo{}, fmt.Errorf("router %s not found", name)
	}
	return router.DeviceInfo(ctx)
}

// ListRoutes proxies route listing to the specified router backend.
func (a *AnyRouterAPI) ListRoutes(ctx context.Context, name string) ([]Route, error) {
	router, ok := a.Router(name)
	if !ok {
		return nil, fmt.Errorf("router %s not found", name)
	}
	return router.ListRoutes(ctx)
}

// AddRoute proxies route creation to the specified router backend.
func (a *AnyRouterAPI) AddRoute(ctx context.Context, name string, route Route) error {
	router, ok := a.Router(name)
	if !ok {
		return fmt.Errorf("router %s not found", name)
	}
	return router.AddRoute(ctx, route)
}

// DeleteRoute removes a route from the specified router backend using the key.
func (a *AnyRouterAPI) DeleteRoute(ctx context.Context, name, key string) error {
	router, ok := a.Router(name)
	if !ok {
		return fmt.Errorf("router %s not found", name)
	}
	return router.DeleteRoute(ctx, key)
}
