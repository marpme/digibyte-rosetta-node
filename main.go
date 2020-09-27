package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/marpme/digibyte-rosetta-node/client"
	"github.com/marpme/digibyte-rosetta-node/configuration"
	"github.com/marpme/digibyte-rosetta-node/provider"
	"github.com/marpme/digibyte-rosetta-node/repository"
	"github.com/marpme/digibyte-rosetta-node/routines"
	"github.com/marpme/digibyte-rosetta-node/services"
)

// NewBlockchainRouter creates a blockchain specific router
// that will handle common routes specified inside the rosetta API specification
func NewBlockchainRouter(cfg *configuration.Config, client client.DigibyteClient, blockRepo *repository.BlockRepository) http.Handler {
	assert, err := asserter.NewServer([]*types.NetworkIdentifier{
		{
			Blockchain:           cfg.NetworkIdentifier.Blockchain,
			Network:              cfg.NetworkIdentifier.Network,
			SubNetworkIdentifier: nil,
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create asserter: %v\n", err)
		os.Exit(1)
	}

	networkAPIController := server.NewNetworkAPIController(services.NewNetworkAPIService(client), assert)
	blockAPIController := server.NewBlockAPIController(services.NewBlockAPIService(client, blockRepo), assert)
	return server.NewRouter(networkAPIController, blockAPIController)
}

func main() {
	configPath := os.Getenv(configuration.ConfigPath)
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := configuration.New(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse config: %v\n", err)
		os.Exit(1)
	}

	rclient := provider.CreateRedisDB(cfg, provider.BlockDB)
	// rAccountClient := provider.CreateRedisDB(cfg, provider.AccountDB)

	blockRepo := repository.NewBlockRepository(rclient)
	// acccountRepo := repository.NewTxRepository(rAccountClient)

	client := client.NewDigibyteClient(cfg)
	router := NewBlockchainRouter(cfg, client, blockRepo)
	var wg sync.WaitGroup
	wg.Add(2)

	go routines.StartRouterCapabilities(&wg, cfg, router)
	go routines.StartSyncingCapabilities(&wg, blockRepo, client)

	wg.Wait()
}
