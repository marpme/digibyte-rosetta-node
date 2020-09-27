package utils

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/coinbase/rosetta-sdk-go/types"
)

func mapTransactions(txs []btcjson.TxRawResult) []*types.TransactionIdentifier {
	var transactionsIdentifiers []*types.TransactionIdentifier
	for _, tx := range txs {
		transactionsIdentifiers = append(transactionsIdentifiers, &types.TransactionIdentifier{
			Hash: tx.Hash,
		})
	}

	return transactionsIdentifiers
}

func MapBlock(block *btcjson.GetBlockVerboseTxResult) *types.BlockResponse {
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

func MapBlockWithTransaction(block *btcjson.GetBlockVerboseTxResult) *types.BlockResponse {
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
