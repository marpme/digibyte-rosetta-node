package client

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/marpme/digibyte-rosetta-node/configuration"
)

// DigibyteClient is the Digibyte blockchain client interface
type DigibyteClient interface {
	// GetBlock returns the digibyte block at given height.
	GetBlock(ctx context.Context, height int64) (*btcjson.GetBlockVerboseTxResult, error)

	// GetBlock returns the Digibyte block with a given hash.
	GetBlockByHash(ctx context.Context, hash string) (*btcjson.GetBlockVerboseTxResult, error)

	// GetBlock returns the Digibyte block with a given hash.
	GetBlockByHashWithTransaction(ctx context.Context, hash string) (*btcjson.GetBlockVerboseTxResult, error)

	// GetLatestBlock returns the latest Digibyte block.
	GetLatestBlock(ctx context.Context) (*btcjson.GetBlockVerboseTxResult, error)

	// GetStatus returns the status overview of the node.
	GetStatus(ctx context.Context) (*btcjson.GetBlockChainInfoResult, error)

	// GetConfig returns the config.
	GetConfig() *configuration.Config
}
