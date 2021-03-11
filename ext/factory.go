package ext

import (
	"github.com/bytepowered/flux"
	"github.com/bytepowered/flux/pkg"
)

var (
	typedFactories = make(map[string]flux.Factory, 16)
)

func RegisterFactory(typeName string, factory flux.Factory) {
	typeName = pkg.MustNotEmpty(typeName, "typeName is empty")
	typedFactories[typeName] = pkg.MustNotNil(factory, "Factory is nil").(flux.Factory)
}

func FactoryByType(typeName string) (flux.Factory, bool) {
	typeName = pkg.MustNotEmpty(typeName, "typeName is empty")
	f, o := typedFactories[typeName]
	return f, o
}
