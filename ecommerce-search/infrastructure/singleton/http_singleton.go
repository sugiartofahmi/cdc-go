package singleton

import "go-service/infrastructure/integrations"

func HttpClientSingleton() *integrations.HttpClient {
	return httpClient
}
