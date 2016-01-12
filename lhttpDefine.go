package lhttp

var (
	HEADER_KEY_PUBLISH     = "publish"
	HEADER_KEY_SUBSCRIBE   = "subscribe"
	HEADER_KEY_UNSUBSCRIBE = "unsubscribe"
	HEADER_KEY_UPSTREAM    = "upstream"
)

var (
	//headers max num not size
	headerMax               = 20
	version                 = "1.0"
	protocolName            = "LHTTP"
	protocolNameWithVersion = "LHTTP/1.0"
	protocolLength          = 9
)

var (
	UPSTREAM_HTTP_METHOD_GET  = "GET"
	UPSTREAM_HTTP_METHOD_POST = "POST"
)

var (
	ProcessorMax = 40
)
