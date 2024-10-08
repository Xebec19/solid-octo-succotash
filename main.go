package main

import (
	"fmt"
	"log/slog"

	"github.com/Xebec19/solid-octo-succotash/internal"
	"github.com/Xebec19/solid-octo-succotash/utils"
)

func main() {

	config, err := utils.LoadConfig(".")

	if err != nil {
		slog.Error(err.Error())
	}

	opensearchClient, err := internal.NewOpensearchClient(config.OPENSEARCH_HOST, config.OPENSEARCH_PORT, config.OPENSEARCH_USERNAME, config.OPENSEARCH_PASSWORD)
	if err != nil {
		slog.Error(err.Error())
	}

	serverConfig := &internal.Server{
		Port:          config.SERVER_ADDRESS,
		OpensearchAPI: opensearchClient,
	}

	srv := serverConfig.NewServer()

	slog.Info(fmt.Sprintf("Starting server on port %s", serverConfig.Port))
	slog.Error(srv.ListenAndServe().Error())
}
