package internal

// 用于隐藏内部实现的Key
type Key struct {
	id string
}

var (
	keyRequestId  = Key{id: "inter.context.request.id/926820fa-7ad8-4444-9080-d690ce31c93a"}
	keyWebContext = Key{id: "inter.context.web.ctx/890b1fa9-93ad-4b44-af24-85bcbfe646b4"}
)