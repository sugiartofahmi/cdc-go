package config

var (
	OpensearchHost                         = Get("OPENSEARCH_HOST", "")
	OpensearchPort                         = Get("OPENSEARCH_PORT", "9200")
	OpensearchUser                         = Get("OPENSEARCH_USER", "admin")
	OpensearchPassword                     = Get("OPENSEARCH_PASSWORD", "")
	OpensearchUseSSL                       = Get("OPENSEARCH_USE_SSL", "false")
	OpensearchMaxIdleConnsPerHost          = stringToInt(Get("OPENSEARCH_MAX_IDLE_CONNS_PER_HOST", "10"))
	OpensearchResponseHeaderTimeoutSeconds = stringToInt(Get("OPENSEARCH_RESPONSE_HEADER_TIMEOUT_SECONDS", "5"))
)
