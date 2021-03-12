package boot

import (
	"fmt"
	"github.com/bytepowered/flux/flux-node"
	"github.com/bytepowered/flux/flux-node/ext"
	"github.com/bytepowered/flux/flux-node/logger"
	"github.com/spf13/viper"
)

const (
	dynConfigKeyDisable = "disable"
	dynConfigKeyTypeId  = "type-id"
)

type AwareConfig struct {
	Id      string
	TypeId  string
	Config  *flux.Configuration
	Factory flux.Factory
}

// 动态加载Filter
func dynamicFilters() ([]AwareConfig, error) {
	out := make([]AwareConfig, 0)
	for id := range viper.GetStringMap("FILTER") {
		v := viper.Sub("FILTER." + id)
		if v == nil || !v.IsSet(dynConfigKeyTypeId) {
			logger.Infow("Filter configuration is empty or without typeId", "typeId", id)
			continue
		}
		config := flux.NewConfigurationOfViper(v)
		typeId := config.GetString(dynConfigKeyTypeId)
		if config.GetBool(dynConfigKeyDisable) {
			logger.Infow("Filter is DISABLED", "typeId", typeId, "id", id)
			continue
		}
		factory, ok := ext.FactoryByType(typeId)
		if !ok {
			return nil, fmt.Errorf("FilterFactory not found, typeId: %s, name: %s", typeId, id)
		}
		out = append(out, AwareConfig{
			Id:      id,
			TypeId:  typeId,
			Factory: factory,
			Config:  config,
		})
	}
	return out, nil
}