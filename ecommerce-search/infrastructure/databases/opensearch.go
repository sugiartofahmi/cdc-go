package databases

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-service/infrastructure/config"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func NewOpenSearchConnection() (*opensearchapi.Client, error) {
	protocol := "http"
	if config.OpensearchUseSSL == "true" {
		protocol = "https"
	}

	address := fmt.Sprintf("%s://%s:%s", protocol, config.OpensearchHost, config.OpensearchPort)

	transport := &http.Transport{
		MaxIdleConnsPerHost:   config.OpensearchMaxIdleConnsPerHost,
		ResponseHeaderTimeout: time.Duration(config.OpensearchResponseHeaderTimeoutSeconds) * time.Second,
	}

	osCfg := opensearch.Config{
		Addresses: []string{address},
		Username:  config.OpensearchUser,
		Password:  config.OpensearchPassword,
		Transport: transport,
	}

	client, err := opensearchapi.NewClient(opensearchapi.Config{Client: osCfg})
	if err != nil {
		return nil, fmt.Errorf("failed to create opensearch client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping opensearch: %w", err)
	}

	return client, nil
}
