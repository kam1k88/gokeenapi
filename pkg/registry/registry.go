package registry

import (
	"fmt"
	"sync"

	"github.com/kam1k88/gokeenapi/pkg/config"
	"github.com/kam1k88/gokeenapi/pkg/goarapi"
)

// Factory creates router backend instance for provided logical name using the
// loaded configuration. It allows plugging different backend implementations.
type Factory func(name string, cfg *config.GokeenapiConfig) (goarapi.RouterAPI, error)

// Registry keeps mapping between backend identifiers and their factories.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// New returns empty registry instance.
func New() *Registry {
	return &Registry{factories: make(map[string]Factory)}
}

// RegisterBackend registers factory under specific backend identifier.
func (r *Registry) RegisterBackend(kind string, factory Factory) {
	if factory == nil || kind == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[kind] = factory
}

// Create instantiates router backend using registered factory.
func (r *Registry) Create(kind, name string, cfg *config.GokeenapiConfig) (goarapi.RouterAPI, error) {
	r.mu.RLock()
	factory, ok := r.factories[kind]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("backend %s is not registered", kind)
	}
	return factory(name, cfg)
}

// Backends returns list of registered backend identifiers.
func (r *Registry) Backends() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.factories))
	for name := range r.factories {
		out = append(out, name)
	}
	return out
}
