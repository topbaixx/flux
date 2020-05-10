package ext

import (
	"github.com/bytepowered/flux"
	"github.com/bytepowered/flux/pkg"
	"sync"
)

const (
	anyHost = "*"
)

var (
	_hostedSelectors = make(map[string][]flux.Selector)
	_rwLock          sync.RWMutex
)

func AddSelector(s flux.Selector) {
	pkg.RequireNotNil(s, "Selector is nil")
	AddHostedSelector(anyHost, s)
}

func AddHostedSelector(host string, s flux.Selector) {
	pkg.RequireNotNil(s, "Selector is nil")
	defer _rwLock.Unlock()
	_rwLock.Lock()
	if l, ok := _hostedSelectors[host]; ok {
		_hostedSelectors[host] = append(l, s)
	} else {
		_hostedSelectors[host] = []flux.Selector{s}
	}
}

func FindSelectors(host string) []flux.Selector {
	defer _rwLock.RUnlock()
	_rwLock.RLock()
	if hosted, ok := _hostedSelectors[host]; ok {
		return _selectors(hosted)
	} else if anyHost != host {
		if any, ok := _hostedSelectors[anyHost]; ok {
			return _selectors(any)
		}
	}
	return nil
}

func _selectors(src []flux.Selector) []flux.Selector {
	out := make([]flux.Selector, len(src))
	copy(out, src)
	return out
}
