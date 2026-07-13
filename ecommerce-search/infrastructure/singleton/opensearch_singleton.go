package singleton

import "github.com/opensearch-project/opensearch-go/v4/opensearchapi"

func OpensearchSingleton() *opensearchapi.Client {
	return opensearchClient
}
