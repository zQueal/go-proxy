package main

import (
	"fmt"
)

type Route interface {
	Start()
	Stop()
}

func NewRoute(cfg *ProxyConfig) (Route, error) {
	if isStreamScheme(cfg.Scheme) {
		id := cfg.GetID()
		if streamRoutes.Contains(id) {
			return nil, fmt.Errorf("duplicated %s stream %s, ignoring", cfg.Scheme, id)
		}
		route, err := NewStreamRoute(cfg)
		if err != nil {
			return nil, err
		}
		streamRoutes.Set(id, route)
		return route, nil
	} else {
		httpRoutes.Ensure(cfg.Alias)
		route, err := NewHTTPRoute(cfg)
		if err != nil {
			return nil, err
		}
		httpRoutes.Get(cfg.Alias).Add(cfg.Path, route)
		return route, nil
	}
}

func isValidScheme(s string) bool {
	for _, v := range ValidSchemes {
		if v == s {
			return true
		}
	}
	return false
}

func isStreamScheme(s string) bool {
	for _, v := range StreamSchemes {
		if v == s {
			return true
		}
	}
	return false
}

// id    -> target
type StreamRoutes = SafeMap[string, StreamRoute]

var streamRoutes = NewSafeMap[string, StreamRoute]()
