package flux

const (
	KeyConfigRootRegistry = "Registry"
	KeyConfigRegistryId   = "registry-id"
)

// Registry Endpoint注册元数据事件监听
type Registry interface {
	// 监听接收元数据中心的配置变化
	Watch() (events chan<- EndpointEvent, err error)
}
