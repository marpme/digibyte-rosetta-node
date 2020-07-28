package utils

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/coinbase/rosetta-sdk-go/types"
)

func mapTransactions(txs []string) []*types.TransactionIdentifier {
	var transactionsIdentifiers []*types.TransactionIdentifier
	for i := 0; i < len(txs); i++ {
		transactionsIdentifiers = append(transactionsIdentifiers, &types.TransactionIdentifier{
			Hash: txs[i],
		})
	}

	return transactionsIdentifiers
}

func MapBlock(block *btcjson.GetBlockVerboseResult) *types.BlockResponse {
	return &types.BlockResponse{
		Block: &types.Block{
			BlockIdentifier: &types.BlockIdentifier{
				Hash: block.Hash,
			},
			ParentBlockIdentifier: &types.BlockIdentifier{
				Hash: block.PreviousHash,
			},
			Timestamp: block.Time,
		},
		OtherTransactions: mapTransactions(block.Tx),
	}
}

func MapBlockWithTransaction(block *btcjson.GetBlockVerboseResult) *types.BlockResponse {
	return &types.BlockResponse{
		Block: &types.Block{
			BlockIdentifier: &types.BlockIdentifier{
				Hash: block.Hash,
			},
			ParentBlockIdentifier: &types.BlockIdentifier{
				Hash: block.PreviousHash,
			},
			Timestamp: block.Time,
		},
		OtherTransactions: mapTransactions(block.Tx),
	}
}
