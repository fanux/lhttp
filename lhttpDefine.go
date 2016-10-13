package lhttp

// CRLF is the end of text line
var CRLF = "\r\n"

var (
	// HEADER_KEY_PUBLISH
	HEADER_KEY_PUBLISH = "publish"
	// HEADER_KEY_SUBSCRIBE
	HEADER_KEY_SUBSCRIBE = "subscribe"
	// HEADER_KEY_UNSUBSCRIBE
	HEADER_KEY_UNSUBSCRIBE = "unsubscribe"
	// HEADER_KEY_UPSTREAM
	HEADER_KEY_UPSTREAM = "upstream"
	// HEADER_KEY_MULTIPART
	HEADER_KEY_MULTIPART = "multipart"
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
	// UPSTREAM_HTTP_METHOD_GET
	UPSTREAM_HTTP_METHOD_GET = "GET"
	// UPSTREAM_HTTP_METHOD_POST
	UPSTREAM_HTTP_METHOD_POST = "POST"
)

var (
	// ProcessorMax
	ProcessorMax = 40
)
