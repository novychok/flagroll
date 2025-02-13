package platformapiv1

import (
	"context"
	"net/http"
	"sync"
)

type Key struct{}

type WebContext struct {
	mu   sync.RWMutex
	Keys map[string]any
}

func (c *WebContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}

	c.Keys[key] = value
}

func (c *WebContext) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.Keys[key]
	return
}

func ContextMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), Key{}, &WebContext{}))
			next.ServeHTTP(w, r)
		})
	}
}

func ContextFromRequest(r *http.Request) *WebContext {
	return r.Context().Value(Key{}).(*WebContext)
}
